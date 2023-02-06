package core

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	merkle "gitlab.com/rarimo/go-merkle"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/operation"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/utils"
	"golang.org/x/exp/slices"
)

type core struct {
	log  *logan.Entry
	core rarimocore.QueryClient
	tm   tokenmanager.QueryClient
}

type Core interface {
	GetTransfers(ctx context.Context, confirmationID string) ([]TransferDetails, error)
	GetTransfer(ctx context.Context, confirmationID string, transferID string) (*TransferDetails, error)
	GetConfirmation(ctx context.Context, confirmationID string) (*rarimocore.Confirmation, error)
}

func NewCore(cfg config.Config) Core {
	return &core{
		core: rarimocore.NewQueryClient(cfg.Cosmos()),
		tm:   tokenmanager.NewQueryClient(cfg.Cosmos()),
		log:  cfg.Log().WithField("service", "core"),
	}
}

func (c *core) GetTransfers(ctx context.Context, confirmationID string) ([]TransferDetails, error) {
	log := c.log.WithField("merkle_root", confirmationID)
	log.Info("processing a confirmation")

	confirmation, err := c.GetConfirmation(ctx, confirmationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch the confirmation")
	}
	networks, err := c.tm.Params(ctx, new(tokenmanager.QueryParamsRequest))
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch the network params")
	}

	transfers := make([]TransferDetails, 0, len(confirmation.Indexes))
	operations := []*operation.TransferContent{}
	contents := []merkle.Content{}

	for _, id := range confirmation.Indexes {
		operation, err := c.core.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: id})
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch the operation")
		}
		transfer := rarimocore.Transfer{}
		if err := transfer.Unmarshal(operation.Operation.Details.Value); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal the transfer")
		}

		token, err := c.tm.Info(ctx, &tokenmanager.QueryGetInfoRequest{Index: transfer.TokenIndex})
		if err != nil {
			return nil, errors.Wrap(err, "error getting token info entry")
		}

		tokenDetails, err := c.tm.Item(ctx, &tokenmanager.QueryGetItemRequest{
			TokenAddress: token.Info.Chains[transfer.ToChain].TokenAddress,
			TokenId:      token.Info.Chains[transfer.ToChain].TokenId,
			Chain:        transfer.ToChain,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error getting token item entry")
		}

		content, err := pkg.GetTransferContent(
			&tokenDetails.Item,
			networks.Params.Networks[transfer.ToChain],
			&transfer,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get transfer content")
		}
		contents = append(contents, content)
		operations = append(operations, content)

		transfers = append(transfers, TransferDetails{
			Transfer:     transfer,
			Token:        token.Info,
			TokenDetails: tokenDetails.Item,
			Signature:    confirmation.SignatureECDSA,
			Origin:       hexutil.Encode(content.Origin[:]),
		})
	}

	tree := merkle.NewTree(crypto.Keccak256, contents...)
	for i, operation := range operations {
		rawPath, ok := tree.Path(operation)
		if !ok {
			panic(fmt.Errorf("failed to build Merkle tree"))
		}
		transfers[i].MerklePath = make([][32]byte, 0, len(rawPath))
		for _, hash := range rawPath {
			transfers[i].MerklePath = append(transfers[i].MerklePath, utils.ToByte32(hash))
		}
	}

	return transfers, nil
}

func (c *core) GetConfirmation(ctx context.Context, confirmationID string) (*rarimocore.Confirmation, error) {
	confirmation, err := c.core.Confirmation(ctx, &rarimocore.QueryGetConfirmationRequest{
		Root: confirmationID,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch the confirmation")
	}
	if confirmation == nil {
		return nil, errors.New("confirmation not found")
	}

	return &confirmation.Confirmation, nil
}

func (c *core) GetTransfer(ctx context.Context, confirmationID string, transferID string) (*TransferDetails, error) {
	transfers, err := c.GetTransfers(ctx, confirmationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transfers")
	}
	transferI := slices.IndexFunc(transfers, func(t TransferDetails) bool {
		return t.Origin == transferID
	})
	if transferI == -1 {
		return nil, errors.New("transfer not found")
	}

	return &transfers[transferI], nil
}
