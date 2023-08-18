package relayer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/core"
	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/data/pg"
	"gitlab.com/rarimo/relayer-svc/pkg/polygonid/contracts"
)

var (
	ErrChainNotFound = errors.New("chain not found")
	ErrEntryNotFound = errors.New("entry not found")
)

type Service struct {
	log     *logan.Entry
	core    *core.Core
	chains  *config.EVM
	storage *pg.Storage
}

func NewService(cfg config.Config) *Service {
	return &Service{
		log:     cfg.Log(),
		core:    core.NewCore(cfg),
		chains:  cfg.EVM(),
		storage: pg.New(cfg.DB()),
	}
}

func (c *Service) Relay(ctx context.Context, state string, chainName string) (string, error) {
	chain, ok := c.chains.GetChainByName(chainName)
	if !ok {
		return "", ErrChainNotFound
	}

	entry, err := c.storage.StateQ().StateByID(state, false)
	if err != nil {
		return "", errors.Wrap(err, "failed to get entry by state")
	}

	if entry == nil {
		return "", ErrEntryNotFound
	}

	return c.processIdentityDefaultTransfer(ctx, chain, entry)
}

func (c *Service) processIdentityDefaultTransfer(ctx context.Context, chain *config.EVMChain, entry *data.State) (string, error) {
	opts := chain.TransactorOpts()

	nonce, err := chain.RPC.PendingNonceAt(context.TODO(), chain.SubmitterAddress)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch a nonce")
	}

	opts.Nonce = big.NewInt(int64(nonce))

	opts.GasPrice, err = chain.RPC.SuggestGasPrice(context.TODO())
	if err != nil {
		return "", errors.Wrap(err, "failed to get suggested gas price")
	}

	details, err := c.core.GetIdentityDefaultTransferProof(ctx, entry.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to get operation proof details")
	}

	replacedState := new(big.Int).SetBytes(hexutil.MustDecode(details.Operation.ReplacedStateHash))
	replacedGIST := new(big.Int).SetBytes(hexutil.MustDecode(details.Operation.ReplacedGISTHash))

	stateInfo, err := getStateInfo(details.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to get state info from transfer")
	}

	gistRootInfo, err := getGistRootInfo(details.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to get gist root info from transfer")
	}

	contract, err := contracts.NewLightweightStateV2(chain.ContractAddress, chain.RPC)
	if err != nil {
		return "", errors.Wrap(err, "failed to create contract instance")
	}

	tx, err := contract.SignedTransitState(opts, replacedState, replacedGIST, stateInfo, gistRootInfo, details.Proof)

	if err != nil {
		return "", errors.Wrap(err, "failed to send state transition tx")
	}

	//go func() {
	//	log := c.log.WithField("state", entry.ID).WithField("operation_id", entry.Operation)
	//
	//	receipt, err := bind.WaitMined(context.Background(), chain.RPC, tx)
	//	if err != nil {
	//		log.WithError(err).Error("failed to wait for state transition tx")
	//		return
	//	}
	//
	//	if receipt.Status == 0 {
	//		log.WithError(err).WithFields(logan.F{
	//			"receipt": utils.Prettify(receipt),
	//			"chain":   chain.Name,
	//		}).Error("failed to wait for state transition tx")
	//		return
	//	}
	//
	//	log.
	//		WithFields(logan.F{
	//			"tx_id":        tx.Hash(),
	//			"tx_index":     receipt.TransactionIndex,
	//			"block_number": receipt.BlockNumber,
	//			"gas_used":     receipt.GasUsed,
	//		}).
	//		Info("evm transaction confirmed")
	//}()

	return tx.Hash().Hex(), nil
}

func getStateInfo(transfer *rarimocore.IdentityDefaultTransfer) (state contracts.ILightweightStateV2StateData, err error) {
	state.Id = new(big.Int).SetBytes(hexutil.MustDecode(transfer.Id))

	state.State = new(big.Int).SetBytes(hexutil.MustDecode(transfer.StateHash))

	state.ReplacedByState = new(big.Int).SetBytes(hexutil.MustDecode(transfer.StateReplacedBy))

	var ok bool
	state.CreatedAtTimestamp, ok = big.NewInt(0).SetString(transfer.StateCreatedAtTimestamp, 10)
	if !ok {
		return contracts.ILightweightStateV2StateData{}, errors.New("failed to decode state created at timestamp")
	}

	state.CreatedAtBlock, ok = big.NewInt(0).SetString(transfer.StateCreatedAtBlock, 10)
	if !ok {
		return contracts.ILightweightStateV2StateData{}, errors.New("failed to decode state created at block")
	}

	return
}

func getGistRootInfo(transfer *rarimocore.IdentityDefaultTransfer) (gistRoot contracts.ILightweightStateV2GistRootData, err error) {
	gistRoot.Root = new(big.Int).SetBytes(hexutil.MustDecode(transfer.GISTHash))

	gistRoot.ReplacedByRoot = new(big.Int).SetBytes(hexutil.MustDecode(transfer.GISTReplacedBy))

	var ok bool
	gistRoot.CreatedAtTimestamp, ok = big.NewInt(0).SetString(transfer.GISTCreatedAtTimestamp, 10)
	if !ok {
		return contracts.ILightweightStateV2GistRootData{}, errors.New("failed to decode GIST created at timestamp")
	}

	gistRoot.CreatedAtBlock, ok = big.NewInt(0).SetString(transfer.GISTCreatedAtBlock, 10)
	if !ok {
		return contracts.ILightweightStateV2GistRootData{}, errors.New("failed to decode GIST created at block")
	}

	return
}
