package relayer

import (
	"context"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"

	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data"
)

const (
	MaxRetries = 3

	prefetchLimit = 10
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 100
)

type relayer struct {
	log   *logan.Entry
	queue rmq.Queue
}

type relayerConsumer struct {
	log          *logan.Entry
	rarimocore   rarimocore.QueryClient
	tokenmanager tokenmanager.QueryClient
	evm          *config.EVM
	solana       *config.Solana
	near         *config.Near
	queue        rmq.Queue
}

func Run(cfg config.Config, ctx context.Context) {
	log := cfg.Log().WithField("service", "relayer")
	r := relayer{
		log:   log,
		queue: cfg.Redis().OpenRelayQueue(),
	}

	if err := r.queue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		panic(errors.Wrap(err, "failed to start consuming the relay queue"))
	}

	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("relay-consumer-%d", i)
		if _, err := r.queue.AddConsumer(name, newConsumer(cfg, name)); err != nil {
			panic(err)
		}
	}

	<-ctx.Done()
	<-r.queue.StopConsuming()
	r.log.Info("finished consuming relayer queue")
}

func newConsumer(cfg config.Config, id string) *relayerConsumer {
	return &relayerConsumer{
		log:          cfg.Log().WithField("service", id),
		rarimocore:   rarimocore.NewQueryClient(cfg.Cosmos()),
		tokenmanager: tokenmanager.NewQueryClient(cfg.Cosmos()),
		evm:          cfg.EVM(),
		solana:       cfg.Solana(),
		near:         cfg.Near(),
		queue:        cfg.Redis().OpenRelayQueue(),
	}
}

func (c *relayerConsumer) Consume(delivery rmq.Delivery) {
	defer func() {
		if err := recover(); err != nil {
			c.log.WithField("err", err).Error("relayer panicked")
		}
	}()

	var task data.RelayTask
	task.Unmarshal(delivery.Payload())

	if err := c.processTransfer(task); err != nil {
		c.log.WithError(err).WithField("transfer_id", task.OperationIndex).Error("failed to process transfer")
		mustReject(delivery)
		c.mustScheduleRetry(task)
		return
	}

	if err := delivery.Ack(); err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to ack the transfer %s", task.OperationIndex)))
	}
}

func (c *relayerConsumer) processTransfer(task data.RelayTask) error {
	c.log.WithField("op_id", task.OperationIndex).Info("processing a transfer")

	operation, err := c.rarimocore.Operation(context.TODO(), &rarimocore.QueryGetOperationRequest{Index: task.OperationIndex})
	if err != nil {
		return errors.Wrap(err, "failed to get the transfer")
	}
	if !operation.Operation.Signed {
		return errors.New("transfer is not signed yet")
	}
	transfer := rarimocore.Transfer{}
	if err := transfer.Unmarshal(operation.Operation.Details.Value); err != nil {
		return errors.Wrap(err, "failed to unmarshal  transfer")
	}

	tokenInfo, err := c.tokenmanager.Info(context.TODO(), &tokenmanager.QueryGetInfoRequest{Index: transfer.TokenIndex})
	if err != nil {
		return errors.Wrap(err, "failed to get token info")
	}
	token, found := tokenInfo.Info.Chains[transfer.ToChain]
	if !found {
		return fmt.Errorf("unknown toChain = %s", transfer.ToChain)
	}
	tokenDetails, err := c.tokenmanager.Item(context.TODO(), &tokenmanager.QueryGetItemRequest{
		TokenAddress: token.TokenAddress,
		TokenId:      token.TokenId,
		Chain:        transfer.ToChain,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get token details")
	}

	network, err := c.tokenmanager.Params(context.TODO(), new(tokenmanager.QueryParamsRequest))
	if err != nil {
		return errors.Wrap(err, "error getting network info entry")
	}

	c.log.
		WithFields(logan.F{
			"op_id":      task.OperationIndex,
			"to":         transfer.Receiver,
			"token_type": tokenDetails.Item.TokenType,
			"to_chain":   transfer.ToChain,
		}).
		Info("relaying a transfer")

	if chain, found := c.evm.GetChainByName(transfer.ToChain); found {
		if err := c.processEVMTransfer(task, transfer, *token, tokenDetails.Item, network.Params, chain); err != nil {
			return errors.Wrap(err, "failed to process ethereum transfer")
		}
	} else if transfer.ToChain == "Solana" {
		if err := c.processSolanaTransfer(task, transfer, *token, tokenDetails.Item, network.Params); err != nil {
			return errors.Wrap(err, "failed to process solana transfer")
		}
	} else if transfer.ToChain == "Near" {
		if err := c.processNearTransfer(task, transfer, *token, tokenDetails.Item, network.Params); err != nil {
			return errors.Wrap(err, "failed to process near transfer")
		}
	} else {
		return fmt.Errorf("toChain = %s is not supported", transfer.ToChain)
	}

	return nil
}

func mustReject(delivery rmq.Delivery) {
	if err := delivery.Reject(); err != nil {
		panic(errors.Wrap(err, "failed to reject the task"))
	}
}

func (c *relayerConsumer) mustScheduleRetry(task data.RelayTask) {
	/**
	TODO:
		- add exponential backoff
		- distinguish retryable and non-retryable errors
		- set up a dead letter queue
	*/
	if task.RetriesLeft == 0 {
		return
	}

	task.RetriesLeft--
	if err := c.queue.PublishBytes(task.Marshal()); err != nil {
		panic(errors.Wrap(err, "failed to schedule the retry"))
	}

}
