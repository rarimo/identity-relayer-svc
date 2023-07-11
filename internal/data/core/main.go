package core

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/accounts/abi"
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

var (
	proofType, _ = abi.NewType("bytes32[]", "", nil)
	sigType, _   = abi.NewType("bytes", "", nil)
	proofArgs    = abi.Arguments{
		{
			Name: "path",
			Type: proofType,
		},
		{
			Name: "signature",
			Type: sigType,
		},
	}
)

type core struct {
	log  *logan.Entry
	core rarimocore.QueryClient
	tm   tokenmanager.QueryClient
}

type Core interface {
	GetIdentityDefaultTransfer(ctx context.Context, confirmationID, operationID string) (*IdentityTransferDetails, error)
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

func (c *core) GetIdentityDefaultTransfers(ctx context.Context, confirmationID string) ([]IdentityTransferDetails, error) {
	confirmation, err := c.GetConfirmation(ctx, confirmationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch the confirmation")
	}

	contents := make([]merkle.Content, len(confirmation.Indexes))
	result := make([]IdentityTransferDetails, len(confirmation.Indexes))

	for i, idx := range confirmation.Indexes {
		rawOp, err := c.core.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: idx})
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch the operation")
		}

		merkleContent, err := getMerkleContentFromIdentityTransferOp(rawOp.Operation)
		if err != nil {
			return nil, errors.Wrap(err, "failed to make the identity transfer merkle content")
		}

		contents[i] = merkleContent
		result[i].OpIndex = rawOp.Operation.Index
		result[i].Signature = confirmation.SignatureECDSA
	}

	tree := merkle.NewTree(crypto.Keccak256, contents...)

	signature := hexutil.MustDecode(confirmation.SignatureECDSA)
	signature[64] += 27

	for i, content := range contents {
		contentPath, ok := tree.Path(content)
		if !ok {
			continue
		}

		pathHashes := make([]common.Hash, 0, len(contentPath))
		for _, p := range contentPath {
			pathHashes = append(pathHashes, common.BytesToHash(p))
		}

		result[i].Proof, err = proofArgs.Pack(contentPath, signature)
		if err != nil {
			return nil, errors.Wrap(err, "failed to pack the proof")
		}
	}

	return result, nil
}

func (c *core) GetIdentityDefaultTransfer(ctx context.Context, confirmationID, operationID string) (*IdentityTransferDetails, error) {
	confirmation, err := c.GetConfirmation(ctx, confirmationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch the confirmation")
	}

	contents := make([]merkle.Content, 0, len(confirmation.Indexes))
	var targetContent merkle.Content

	result := IdentityTransferDetails{
		Signature: confirmation.SignatureECDSA,
	}

	for _, idx := range confirmation.Indexes {
		rawOp, err := c.core.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: idx})
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch the operation")
		}

		if rawOp.Operation.OperationType != rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER {
			continue
		}

		merkleContent, err := getMerkleContentFromIdentityTransferOp(rawOp.Operation)
		if err != nil {
			return nil, errors.Wrap(err, "failed to make the identity transfer merkle content")
		}

		contents = append(contents, merkleContent)

		if rawOp.Operation.Index == operationID {
			targetContent = merkleContent
			result.OpIndex = rawOp.Operation.Index
		}
	}

	tree := merkle.NewTree(crypto.Keccak256, contents...)

	path, _ := tree.Path(targetContent)

	pathHashes := make([]common.Hash, 0, len(path))
	result.MerklePath = make([][32]byte, 0, len(path))
	for _, p := range path {
		pathHashes = append(pathHashes, common.BytesToHash(p))
		result.MerklePath = append(result.MerklePath, utils.ToByte32(p))
	}

	signature := hexutil.MustDecode(confirmation.SignatureECDSA)
	signature[64] += 27

	result.Proof, err = proofArgs.Pack(pathHashes, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode the proof")
	}

	return &result, nil
}

func getMerkleContentFromIdentityTransferOp(operation rarimocore.Operation) (merkle.Content, error) {
	transfer, err := pkg.GetIdentityDefaultTransfer(operation)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch the transfer")
	}

	return pkg.GetIdentityDefaultTransferContent(transfer)
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

		item, err := c.tm.ItemByOnChainItem(ctx, &tokenmanager.QueryGetItemByOnChainItemRequest{
			Chain:   transfer.From.Chain,
			Address: transfer.From.Address,
			TokenID: transfer.From.TokenID,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error getting item entry")
		}

		collection, err := c.tm.Collection(ctx, &tokenmanager.QueryGetCollectionRequest{
			Index: item.Item.Collection,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error getting collection entry")
		}

		collectionData, err := c.tm.CollectionData(ctx, &tokenmanager.QueryGetCollectionDataRequest{
			Chain:   transfer.From.Chain,
			Address: transfer.From.Address,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error getting collection data entry")
		}

		var netParams *tokenmanager.NetworkParams
		for _, net := range networks.Params.Networks {
			if net.Name == transfer.From.Chain {
				netParams = net
				break
			}
		}

		if netParams == nil {
			return nil, errors.Wrap(err, "no network params found")
		}

		content, err := pkg.GetTransferContent(
			collection.Collection,
			collectionData.Data,
			item.Item,
			*netParams,
			&transfer,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get transfer content")
		}
		contents = append(contents, content)
		operations = append(operations, content)

		transfers = append(transfers, TransferDetails{
			Transfer:      transfer,
			DstCollection: collectionData.Data,
			Item:          item.Item,
			Signature:     confirmation.SignatureECDSA,
			Origin:        hexutil.Encode(content.Origin[:]),
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
