package relayer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/utils"
	"gitlab.com/rarimo/relayer-svc/pkg/polygonid/contracts"
)

const (
	MaxRetries = 0

	prefetchLimit = 10
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 100
)

type StateTransitioner interface {
	SignedTransitState(opts *bind.TransactOpts, prevState_ *big.Int, stateInfo_ contracts.IStateStateInfo, gistRootInfo_ contracts.IStateGistRootInfo, signature_ []byte) (*evmtypes.Transaction, error)
}

type relayer struct {
	log   *logan.Entry
	queue rmq.Queue
}

type relayerConsumer struct {
	log               *logan.Entry
	rarimocore        rarimocore.QueryClient
	targetChain       *config.EVMChain
	stateTransitioner StateTransitioner
	queue             rmq.Queue
}

func Run(cfg config.Config, ctx context.Context) {
	log := cfg.Log().WithField("service", "relayer")
	r := relayer{
		log:   log,
		queue: cfg.Redis().OpenRelayQueue(),
	}

	targetChain, ok := cfg.EVM().GetChainByName(cfg.Relay().TargetChain)
	if !ok {
		panic(fmt.Errorf("unknown target chain [%s]", cfg.Relay().TargetChain))
	}

	if err := r.queue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		panic(errors.Wrap(err, "failed to start consuming the relay queue"))
	}

	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("relay-consumer-%d", i)
		if _, err := r.queue.AddConsumer(name, newConsumer(cfg, name, targetChain)); err != nil {
			panic(err)
		}
	}

	<-ctx.Done()
	<-r.queue.StopConsuming()
	r.log.Info("finished consuming relayer queue")
}

func newConsumer(cfg config.Config, id string, chain *config.EVMChain) *relayerConsumer {
	return &relayerConsumer{
		log:         cfg.Log().WithField("service", id),
		rarimocore:  rarimocore.NewQueryClient(cfg.Cosmos()),
		targetChain: chain,
		queue:       cfg.Redis().OpenRelayQueue(),
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

	if err := c.processIdentityTransfer(task); err != nil {
		c.log.WithError(err).WithField("transfer_id", task.OperationIndex).Error("failed to process transfer")
		mustReject(delivery)
		c.mustScheduleRetry(task)
		return
	}

	if err := delivery.Ack(); err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to ack the transfer %s", task.OperationIndex)))
	}
}

func (c *relayerConsumer) processIdentityTransfer(task data.RelayTask) error {
	log := c.log.WithField("op_id", task.OperationIndex)

	log.Info("processing operation")
	operation, err := c.rarimocore.Operation(context.TODO(), &rarimocore.QueryGetOperationRequest{Index: task.OperationIndex})
	if err != nil {
		return errors.Wrap(err, "failed to get the operation")
	}
	if operation.Operation.Status != rarimocore.OpStatus_SIGNED {
		return errors.New("operation is not signed yet")
	}

	if operation.Operation.OperationType != rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER {
		return fmt.Errorf("unknown operation type %s", operation.Operation.OperationType)
	}

	return c.processIdentityDefaultTransfer(hexutil.MustDecode(task.Proof), operation.Operation.Details.Value)
}

func (c *relayerConsumer) processIdentityDefaultTransfer(proof []byte, raw []byte) error {
	var transfer rarimocore.IdentityDefaultTransfer
	if err := transfer.Unmarshal(raw); err != nil {
		return errors.Wrap(err, "failed to unmarshal identity transfer")
	}

	opts := c.targetChain.TransactorOpts()
	nonce, err := c.targetChain.RPC.PendingNonceAt(context.TODO(), c.targetChain.SubmitterAddress)
	if err != nil {
		return errors.Wrap(err, "failed to fetch a nonce")
	}
	opts.Nonce = big.NewInt(int64(nonce))
	gasPrice, err := c.targetChain.RPC.SuggestGasPrice(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to get suggested gas price")
	}
	opts.GasPrice = gasPrice
	opts.GasLimit = uint64(300000)

	replacedState := new(big.Int).SetBytes(hexutil.MustDecode(transfer.ReplacedStateHash))

	stateInfo, err := getStateInfo(transfer)
	if err != nil {
		return errors.Wrap(err, "failed to get state info from transfer")
	}

	gistRootInfo, err := getGistRootInfo(transfer)
	if err != nil {
		return errors.Wrap(err, "failed to get gist root info from transfer")
	}

	tx, err := c.stateTransitioner.SignedTransitState(opts, replacedState, stateInfo, gistRootInfo, proof)

	if err != nil {
		return errors.Wrap(err, "failed to send state transition tx")
	}

	c.log.WithField("tx_hash", tx.Hash().Hex()).Debug("state transition tx sent")

	receipt, err := bind.WaitMined(context.Background(), c.targetChain.RPC, tx)
	if err != nil {
		return errors.Wrap(err, "failed to wait for state transition tx")
	}

	if err != nil {
		return errors.Wrap(err, "failed to wait for the transaction to be mined")
	}
	if receipt.Status == 0 {
		return errors.Wrap(err, "transaction failed", logan.F{
			"receipt": utils.Prettify(receipt),
			"chain":   transfer.Chain,
		})
	}

	c.log.
		WithFields(logan.F{
			"tx_id":        tx.Hash(),
			"tx_index":     receipt.TransactionIndex,
			"block_number": receipt.BlockNumber,
			"gas_used":     receipt.GasUsed,
		}).
		Info("evm transaction confirmed")

	return nil
}

func getStateInfo(transfer rarimocore.IdentityDefaultTransfer) (state contracts.IStateStateInfo, err error) {
	state.Id = new(big.Int).SetBytes(hexutil.MustDecode(transfer.Id))

	state.State = new(big.Int).SetBytes(hexutil.MustDecode(transfer.StateHash))

	state.ReplacedByState = new(big.Int).SetBytes(hexutil.MustDecode(transfer.StateReplacedBy))

	var ok bool
	state.CreatedAtTimestamp, ok = big.NewInt(0).SetString(transfer.StateCreatedAtTimestamp, 10)
	if !ok {
		return contracts.IStateStateInfo{}, errors.New("failed to decode state created at timestamp")
	}

	state.ReplacedAtTimestamp, ok = big.NewInt(0).SetString(transfer.StateReplacedAtTimestamp, 10)
	if !ok {
		return contracts.IStateStateInfo{}, errors.New("failed to decode state replaced at timestamp")
	}

	state.CreatedAtBlock, ok = big.NewInt(0).SetString(transfer.StateCreatedAtBlock, 10)
	if !ok {
		return contracts.IStateStateInfo{}, errors.New("failed to decode state created at block")
	}

	state.ReplacedAtBlock, ok = big.NewInt(0).SetString(transfer.StateReplacedAtBlock, 10)
	if !ok {
		return contracts.IStateStateInfo{}, errors.New("failed to decode state replaced at block")
	}

	return
}

func getGistRootInfo(transfer rarimocore.IdentityDefaultTransfer) (gistRoot contracts.IStateGistRootInfo, err error) {
	gistRoot.Root = new(big.Int).SetBytes(hexutil.MustDecode(transfer.GISTHash))

	gistRoot.ReplacedByRoot = new(big.Int).SetBytes(hexutil.MustDecode(transfer.GISTReplacedBy))

	var ok bool
	gistRoot.CreatedAtTimestamp, ok = big.NewInt(0).SetString(transfer.GISTCreatedAtTimestamp, 10)
	if !ok {
		return contracts.IStateGistRootInfo{}, errors.New("failed to decode GIST created at timestamp")
	}

	gistRoot.ReplacedAtTimestamp, ok = big.NewInt(0).SetString(transfer.GISTReplacedAtTimestamp, 10)
	if !ok {
		return contracts.IStateGistRootInfo{}, errors.New("failed to decode GIST replaced at timestamp")
	}

	gistRoot.CreatedAtBlock, ok = big.NewInt(0).SetString(transfer.GISTCreatedAtBlock, 10)
	if !ok {
		return contracts.IStateGistRootInfo{}, errors.New("failed to decode GIST created at block")
	}

	gistRoot.ReplacedAtBlock, ok = big.NewInt(0).SetString(transfer.GISTReplacedAtBlock, 10)
	if !ok {
		return contracts.IStateGistRootInfo{}, errors.New("failed to decode GIST replaced at block")
	}

	return
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
