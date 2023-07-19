package core

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	merkle "gitlab.com/rarimo/go-merkle"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
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
}

func NewCore(cfg config.Config) Core {
	return &core{
		core: rarimocore.NewQueryClient(cfg.Cosmos()),
		tm:   tokenmanager.NewQueryClient(cfg.Cosmos()),
		log:  cfg.Log().WithField("service", "core"),
	}
}

func (c *core) GetIdentityDefaultTransfer(ctx context.Context, confirmationID, operationID string) (*IdentityTransferDetails, error) {
	c.log.Debugf("Starting proof generation for operation: %s", operationID)

	confirmation, err := c.GetConfirmation(ctx, confirmationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch the confirmation")
	}

	var targetContent merkle.Content
	var result IdentityTransferDetails

	operations := make([]rarimocore.Operation, 0, len(confirmation.Indexes))
	for _, idx := range confirmation.Indexes {
		rawOp, err := c.core.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: idx})
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch the operation")
		}

		if rawOp.Operation.Index == operationID && rawOp.Operation.OperationType == rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER {
			targetContent, err = c.getIdentityDefaultTransferContent(rawOp.Operation)
			if err != nil {
				return nil, errors.Wrap(err, "failed to make the identity transfer merkle content", logan.F{
					"operation": rawOp.Operation.Index,
				})
			}

			result.OpIndex = rawOp.Operation.Index
		}

		operations = append(operations, rawOp.Operation)
	}

	contents, err := c.getContents(operations...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get contents")
	}

	tree := merkle.NewTree(crypto.Keccak256, contents...)

	c.log.Debugf("Reconstructed tree root: %s", hexutil.Encode(tree.Root()))

	path, _ := tree.Path(targetContent)

	c.log.Debugf("Reconstructed path length: %d", len(path))

	pathHashes := make([]common.Hash, 0, len(path))
	for _, p := range path {
		pathHashes = append(pathHashes, common.BytesToHash(p))
	}

	signature := hexutil.MustDecode(confirmation.SignatureECDSA)
	signature[64] += 27

	result.Proof, err = proofArgs.Pack(pathHashes, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode the proof")
	}

	c.log.Debugf("Generated proof: %s", hexutil.Encode(result.Proof))

	return &result, nil
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
