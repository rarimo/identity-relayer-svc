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
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/types"

	"gitlab.com/rarimo/relayer-svc/internal/services/bridger"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger/bridge"
)

const (
	MaxRetries = 0

	prefetchLimit = 10
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 100
)

type relayer struct {
	log   *logan.Entry
	queue rmq.Queue
}

type relayerConsumer struct {
	log             *logan.Entry
	rarimocore      rarimocore.QueryClient
	tokenmanager    tokenmanager.QueryClient
	evm             *config.EVM
	bridgerProvider bridger.BridgerProvider
	solana          *config.Solana
	near            *config.Near
	queue           rmq.Queue
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
		log:             cfg.Log().WithField("service", id),
		rarimocore:      rarimocore.NewQueryClient(cfg.Cosmos()),
		tokenmanager:    tokenmanager.NewQueryClient(cfg.Cosmos()),
		evm:             cfg.EVM(),
		solana:          cfg.Solana(),
		near:            cfg.Near(),
		queue:           cfg.Redis().OpenRelayQueue(),
		bridgerProvider: bridger.NewBridgerProvider(cfg),
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

	if err := c.processOperation(task); err != nil {
		c.log.WithError(err).WithField("transfer_id", task.OperationIndex).Error("failed to process transfer")
		mustReject(delivery)
		c.mustScheduleRetry(task)
		return
	}

	if err := delivery.Ack(); err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to ack the transfer %s", task.OperationIndex)))
	}
}

func (c *relayerConsumer) processOperation(task data.RelayTask) error {
	log := c.log.WithField("op_id", task.OperationIndex)

	log.Info("processing operation")
	operation, err := c.rarimocore.Operation(context.TODO(), &rarimocore.QueryGetOperationRequest{Index: task.OperationIndex})
	if err != nil {
		return errors.Wrap(err, "failed to get the operation")
	}
	if operation.Operation.Status != rarimocore.OpStatus_SIGNED {
		return errors.New("operation is not signed yet")
	}

	switch operation.Operation.OperationType {
	case rarimocore.OpType_TRANSFER:
		return c.processTransfer(task, operation.Operation.Details.Value)
	//case rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER:
	//	return c.processIdentityDefaultTransfer(task, operation.Operation.Details.Value)
	default:
		return fmt.Errorf("unknown operation type %s", operation.Operation.OperationType)
	}
}

func (c *relayerConsumer) processTransfer(task data.RelayTask, raw []byte) error {
	transfer := rarimocore.Transfer{}
	if err := transfer.Unmarshal(raw); err != nil {
		return errors.Wrap(err, "failed to unmarshal  transfer")
	}

	dstItem, err := c.tokenmanager.ItemByOnChainItem(context.TODO(),
		&tokenmanager.QueryGetItemByOnChainItemRequest{
			Chain:   transfer.To.Chain,
			Address: transfer.To.Address,
			TokenID: transfer.To.TokenID,
		})
	if err != nil {
		return errors.Wrap(err, "failed to get dst item info")
	}

	collectionData, err := c.tokenmanager.CollectionData(context.TODO(),
		&tokenmanager.QueryGetCollectionDataRequest{
			Chain:   transfer.To.Chain,
			Address: transfer.To.Address,
		})
	if err != nil {
		return errors.Wrap(err, "failed to get collection data")
	}

	collection, err := c.tokenmanager.Collection(context.TODO(),
		&tokenmanager.QueryGetCollectionRequest{
			Index: collectionData.Data.Collection,
		})
	if err != nil {
		return errors.Wrap(err, "failed to get collection")
	}

	transferDetails := core.TransferDetails{
		Transfer:      transfer,
		DstCollection: collectionData.Data,
		Item:          dstItem.Item,
		Signature:     task.Signature,
		Origin:        task.Origin,
		MerklePath:    task.MustParseMerklePath(),
	}

	c.log.
		WithFields(logan.F{
			"to":         transfer.Receiver,
			"token_type": collectionData.Data.TokenType,
			"to_chain":   transfer.To.Chain,
		}).
		Info("relaying a transfer")

	switch {
	case transfer.To.Chain == types.Near:
		err = c.processNearTransfer(task, transfer, dstItem.Item, collectionData.Data, collection.Collection.Meta)
	default:
		bridger := c.bridgerProvider.GetBridger(transfer.To.Chain)
		err = bridger.Withdraw(context.Background(), transferDetails)
	}

	if errors.Cause(err) == bridge.ErrAlreadyWithdrawn {
		c.log.Info("transfer has already been withdrawn")
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to process a transfer")
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
