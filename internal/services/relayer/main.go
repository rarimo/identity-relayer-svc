package relayer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rarimo/identity-relayer-svc/internal/config"
	"github.com/rarimo/identity-relayer-svc/internal/core"
	"github.com/rarimo/identity-relayer-svc/internal/data"
	"github.com/rarimo/identity-relayer-svc/internal/data/pg"
	"github.com/rarimo/identity-relayer-svc/internal/utils"
	"github.com/rarimo/identity-relayer-svc/pkg/polygonid/contracts"
	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrChainNotFound    = errors.New("chain not found")
	ErrEntryNotFound    = errors.New("entry not found")
	ErrAlreadySubmitted = errors.New("already transited")
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

func (c *Service) StateRelays(ctx context.Context, state string) ([]data.Transition, error) {
	entry, err := c.storage.StateQ().StateByIDCtx(ctx, state, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entry by state")
	}

	if entry == nil {
		return nil, ErrEntryNotFound
	}

	transitions, err := c.storage.TransitionQ().TransitionsByStateCtx(ctx, state, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transition")
	}

	return transitions, nil
}

func (c *Service) GistRelays(ctx context.Context, gist string) ([]data.GistTransition, error) {
	entry, err := c.storage.GistQ().GistByIDCtx(ctx, gist, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entry by state")
	}

	if entry == nil {
		return nil, ErrEntryNotFound
	}

	transitions, err := c.storage.GistTransitionQ().GistTransitionsByGistCtx(ctx, gist, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transition")
	}

	return transitions, nil
}

func (c *Service) StateRelay(ctx context.Context, state string, chainName string, waitTxConfirm bool) (txhash string, err error) {
	chain, ok := c.chains.GetChainByName(chainName)
	if !ok {
		return "", ErrChainNotFound
	}

	entry, err := c.storage.StateQ().StateByIDCtx(ctx, state, false)
	if err != nil {
		return "", errors.Wrap(err, "failed to get entry by state")
	}

	if entry == nil {
		return "", ErrEntryNotFound
	}

	if err := c.checkTransitionNotExist(ctx, state, chainName); err != nil {
		return "", err
	}

	return c.processIdentityStateTransfer(ctx, chain, entry, waitTxConfirm)
}

func (c *Service) GistRelay(ctx context.Context, gist string, chainName string, waitTxConfirm bool) (txhash string, err error) {
	chain, ok := c.chains.GetChainByName(chainName)
	if !ok {
		return "", ErrChainNotFound
	}

	entry, err := c.storage.GistQ().GistByIDCtx(ctx, gist, false)
	if err != nil {
		return "", errors.Wrap(err, "failed to get entry by gist")
	}

	if entry == nil {
		return "", ErrEntryNotFound
	}

	if err := c.checkGISTTransitionNotExist(ctx, gist, chainName); err != nil {
		return "", err
	}

	return c.processIdentityGISTTransfer(ctx, chain, entry, waitTxConfirm)
}

func (c *Service) checkTransitionNotExist(ctx context.Context, state, chain string) error {
	transitions, err := c.storage.TransitionQ().TransitionsByStateCtx(ctx, state, false)
	if err != nil {
		return errors.Wrap(err, "failed to get transition")
	}

	if len(transitions) == 0 {
		return nil
	}

	for _, transition := range transitions {
		if transition.Chain == chain {
			return ErrAlreadySubmitted
		}
	}

	return nil
}

func (c *Service) checkGISTTransitionNotExist(ctx context.Context, state, chain string) error {
	transitions, err := c.storage.GistTransitionQ().GistTransitionsByGistCtx(ctx, state, false)
	if err != nil {
		return errors.Wrap(err, "failed to get transition")
	}

	if len(transitions) == 0 {
		return nil
	}

	for _, transition := range transitions {
		if transition.Chain == chain {
			return ErrAlreadySubmitted
		}
	}

	return nil
}

func (c *Service) processIdentityStateTransfer(ctx context.Context, chain *config.EVMChain, entry *data.State, waitTxConfirm bool) (txhash string, err error) {
	opts := chain.TransactorOpts()

	nonce, err := chain.RPC.PendingNonceAt(context.TODO(), chain.SubmitterAddress)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch a nonce")
	}

	opts.Nonce = big.NewInt(int64(nonce))

	opts.GasPrice, err = chain.RPC.SuggestGasPrice(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get suggested gas price")
	}

	details, err := c.core.GetIdentityStateTransferProof(ctx, entry.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to get operation proof details")
	}

	stateInfo, err := getStateInfo(details.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to get state info from transfer")
	}

	contract, err := contracts.NewLightweightStateV2(chain.ContractAddress, chain.RPC)
	if err != nil {
		return "", errors.Wrap(err, "failed to create contract instance")
	}

	prevState := new(big.Int).SetBytes(hexutil.MustDecode(details.Operation.ReplacedStateHash))

	tx, err := contract.SignedTransitStateData(opts, prevState, stateInfo, details.Proof)
	if err != nil {
		c.log.Debugf(
			"Tx args: %s, %v, %s",
			prevState.String(),
			stateInfo,
			hexutil.Encode(details.Proof),
		)
		return "", errors.Wrap(err, "failed to send state transition tx")
	}

	transition := data.Transition{
		Tx:    tx.Hash().String(),
		State: entry.ID,
		Chain: chain.Name,
	}

	if err := c.storage.TransitionQ().Insert(&transition); err != nil {
		c.log.WithError(err).Error("failed to create transition entry")
	}

	if waitTxConfirm {
		c.waitTxConfirmation(ctx, chain, tx)
	}

	return tx.Hash().Hex(), nil
}

func (c *Service) processIdentityGISTTransfer(ctx context.Context, chain *config.EVMChain, entry *data.Gist, waitTxConfirm bool) (txhash string, err error) {
	opts := chain.TransactorOpts()

	nonce, err := chain.RPC.PendingNonceAt(context.TODO(), chain.SubmitterAddress)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch a nonce")
	}

	opts.Nonce = big.NewInt(int64(nonce))

	opts.GasPrice, err = chain.RPC.SuggestGasPrice(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get suggested gas price")
	}

	details, err := c.core.GetIdentityGISTTransferProof(ctx, entry.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to get operation proof details")
	}

	gistInfo, err := getGistRootInfo(details.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to get gist info from transfer")
	}

	contract, err := contracts.NewLightweightStateV2(chain.ContractAddress, chain.RPC)
	if err != nil {
		return "", errors.Wrap(err, "failed to create contract instance")
	}

	prevGist := new(big.Int).SetBytes(hexutil.MustDecode(details.Operation.ReplacedGISTHash))

	tx, err := contract.SignedTransitGISTData(opts, prevGist, gistInfo, details.Proof)
	if err != nil {
		c.log.Debugf(
			"Tx args: %s, %v, %s",
			prevGist.String(),
			gistInfo,
			hexutil.Encode(details.Proof),
		)
		return "", errors.Wrap(err, "failed to send gist transition tx")
	}

	transition := data.GistTransition{
		Tx:    tx.Hash().String(),
		Gist:  entry.ID,
		Chain: chain.Name,
	}

	if err := c.storage.GistTransitionQ().InsertCtx(ctx, &transition); err != nil {
		c.log.WithError(err).Error("failed to create transition entry")
	}

	if waitTxConfirm {
		c.waitTxConfirmation(ctx, chain, tx)
	}

	return tx.Hash().Hex(), nil
}

func (c *Service) waitTxConfirmation(ctx context.Context, chain *config.EVMChain, tx *types.Transaction) {
	receipt, err := bind.WaitMined(ctx, chain.RPC, tx)
	if err != nil {
		c.log.WithError(err).Error("failed to wait for state transition tx")
		return
	}

	if receipt.Status == 0 {
		c.log.WithError(err).WithFields(logan.F{
			"receipt": utils.Prettify(receipt),
			"chain":   chain.Name,
		}).Error("failed to wait for state transition tx")
		return
	}

	c.log.
		WithFields(logan.F{
			"tx_id":        tx.Hash(),
			"tx_index":     receipt.TransactionIndex,
			"block_number": receipt.BlockNumber,
			"gas_used":     receipt.GasUsed,
		}).
		Info("evm transaction confirmed")
}

func getStateInfo(transfer *rarimocore.IdentityStateTransfer) (state contracts.ILightweightStateV2StateData, err error) {
	state.Id = new(big.Int).SetBytes(hexutil.MustDecode(transfer.Id))

	state.State = new(big.Int).SetBytes(hexutil.MustDecode(transfer.StateHash))

	var ok bool
	state.CreatedAtTimestamp, ok = big.NewInt(0).SetString(transfer.StateCreatedAtTimestamp, 10)
	if !ok {
		return contracts.ILightweightStateV2StateData{}, errors.New("failed to decode state created at timestamp")
	}

	state.CreatedAtBlock, ok = big.NewInt(0).SetString(transfer.StateCreatedAtBlock, 10)
	if !ok {
		return contracts.ILightweightStateV2StateData{}, errors.New("failed to decode state created at block")
	}

	// Will not be used in proof
	state.ReplacedByState = big.NewInt(0)
	return
}

func getGistRootInfo(transfer *rarimocore.IdentityGISTTransfer) (gistRoot contracts.ILightweightStateV2GistRootData, err error) {
	gistRoot.Root = new(big.Int).SetBytes(hexutil.MustDecode(transfer.GISTHash))

	var ok bool
	gistRoot.CreatedAtTimestamp, ok = big.NewInt(0).SetString(transfer.GISTCreatedAtTimestamp, 10)
	if !ok {
		return contracts.ILightweightStateV2GistRootData{}, errors.New("failed to decode GIST created at timestamp")
	}

	gistRoot.CreatedAtBlock, ok = big.NewInt(0).SetString(transfer.GISTCreatedAtBlock, 10)
	if !ok {
		return contracts.ILightweightStateV2GistRootData{}, errors.New("failed to decode GIST created at block")
	}

	// Will not be used in proof
	gistRoot.ReplacedByRoot = big.NewInt(0)
	return
}
