package core

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

type Core struct {
	log  *logan.Entry
	core rarimocore.QueryClient
	tm   tokenmanager.QueryClient
}

func NewCore(cfg config.Config) *Core {
	return &Core{
		core: rarimocore.NewQueryClient(cfg.Cosmos()),
		tm:   tokenmanager.NewQueryClient(cfg.Cosmos()),
		log:  cfg.Log().WithField("service", "core"),
	}
}

func (c *Core) GetIdentityDefaultTransferProof(ctx context.Context, operationID string) (*IdentityTransferDetails, error) {
	proof, err := c.core.OperationProof(ctx, &rarimocore.QueryGetOperationProofRequest{Index: operationID})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the operation proof")
	}

	pathHashes := make([]common.Hash, 0, len(proof.Path))
	for _, p := range proof.Path {
		pathHashes = append(pathHashes, common.HexToHash(p))
	}

	signature := hexutil.MustDecode(proof.Signature)
	signature[64] += 27

	operation, err := c.core.Operation(context.TODO(), &rarimocore.QueryGetOperationRequest{Index: operationID})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the operation")
	}

	transfer, err := pkg.GetIdentityDefaultTransfer(operation.Operation)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse operation details")
	}

	result := IdentityTransferDetails{Operation: transfer}

	result.Proof, err = proofArgs.Pack(pathHashes, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode the proof")
	}

	return &result, nil
}
