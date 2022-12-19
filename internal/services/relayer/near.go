package relayer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarify-protocol/near-bridge-go/pkg/client"
	"gitlab.com/rarify-protocol/near-bridge-go/pkg/types"
	"gitlab.com/rarify-protocol/near-bridge-go/pkg/types/action"
	"gitlab.com/rarify-protocol/near-bridge-go/pkg/types/action/base"
	rarimocore "gitlab.com/rarify-protocol/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarify-protocol/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarify-protocol/relayer-svc/internal/data"
	"lukechampine.com/uint128"
)

var nearGasLimit uint64 = 300000000000000

func (c *relayerConsumer) processNearTransfer(
	task data.RelayTask,
	transfer rarimocore.Transfer,
	token tokenmanager.ChainInfo,
	tokenDetails tokenmanager.Item,
	network tokenmanager.Params,
) error {
	log := c.log.WithField("op_id", task.OperationIndex)

	amount, err := parseNearAmount(transfer.Amount)
	if err != nil {
		return errors.Wrap(err, "failed to parse amount")
	}
	rawSignature := hexutil.MustDecode(task.Signature)
	signature := hexutil.Encode(rawSignature[:64])
	withdrawArgs := action.WithdrawArgs{
		Chain:      transfer.ToChain,
		ReceiverID: string(hexutil.MustDecode(transfer.Receiver)),
		Origin:     transfer.Origin,
		Path:       task.MustParseMerklePath(),
		Signatures: []string{signature},
		RecoveryID: rawSignature[64],
	}
	deposit := types.OneYocto

	var act base.Action

	switch tokenDetails.TokenType {
	case tokenmanager.Type_NATIVE:
		args := action.NativeWithdrawArgs{
			Amount:       amount,
			WithdrawArgs: withdrawArgs,
		}
		act = action.NewNativeWithdrawCall(args, nearGasLimit, deposit)
	case tokenmanager.Type_NEAR_FT:
		args := action.FtWithdrawArgs{
			Token:        string(hexutil.MustDecode(token.TokenAddress)),
			Amount:       amount,
			IsWrapped:    tokenDetails.Wrapped,
			WithdrawArgs: withdrawArgs,
		}
		act = action.NewFtWithdrawCall(args, nearGasLimit, deposit)
	case tokenmanager.Type_NEAR_NFT:
		args := action.NftWithdrawArgs{
			Token:        string(hexutil.MustDecode(token.TokenAddress)),
			TokenID:      string(hexutil.MustDecode(token.TokenId)),
			IsWrapped:    tokenDetails.Wrapped,
			WithdrawArgs: withdrawArgs,
		}
		if tokenDetails.Wrapped {
			args.TokenMetadata = &types.NftMetadataView{
				// TODO: fetch the rest of the fields from Horizon
				Title:     tokenDetails.Name,
				Media:     tokenDetails.ImageUri,
				MediaHash: hexutil.MustDecode(tokenDetails.ImageHash),
			}
			deposit = types.NEARToYocto(1)
		}

		act = action.NewNftWithdrawCall(args, nearGasLimit, deposit)
	default:
		return errors.Errorf("invalid solana token type: %d", tokenDetails.TokenType)
	}

	withdrawResp, err := c.near.RPC.TransactionSendAwait(
		client.ContextWithKeyPair(context.TODO(), c.near.SubmitterPrivateKey),
		c.near.SubmitterAddress,
		c.near.BridgeAddress,
		[]base.Action{act},
		client.WithLatestBlock(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to submit a Near transaction")
	}
	if len(withdrawResp.Status.Failure) != 0 {
		log.
			WithField("tx_id", withdrawResp.Transaction.Hash).
			WithField("status_failure", prettify(withdrawResp.Status.Failure)).
			Info("near transaction failed")

		return errors.New("near transaction failed")
	}

	log.WithField("tx_id", withdrawResp.Transaction.Hash).Info("successfully submitted Near transaction")

	return nil
}

func parseNearAmount(raw string) (types.Balance, error) {
	bigAmount, err := getAmountOrDefault(raw, big.NewInt(1))
	if err != nil {
		return types.Balance{}, errors.Wrap(err, "failed to parse amount")
	}

	return types.Balance(uint128.FromBig(bigAmount)), nil
}
