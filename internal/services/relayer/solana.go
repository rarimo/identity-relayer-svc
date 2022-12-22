package relayer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/olegfomenko/solana-go"
	"github.com/olegfomenko/solana-go/rpc"
	confirm "github.com/olegfomenko/solana-go/rpc/sendAndConfirmTransaction"
	"gitlab.com/distributed_lab/logan/v3/errors"

	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"

	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/helpers"
	"gitlab.com/rarimo/solana-program-go/contract"
)

func (c *relayerConsumer) processSolanaTransfer(
	task data.RelayTask,
	transfer rarimocore.Transfer,
	token tokenmanager.ChainInfo,
	tokenDetails tokenmanager.Item,
	network tokenmanager.Params,
) error {
	log := c.log.WithField("op_id", task.OperationIndex)

	receiver := hexutil.MustDecode(transfer.Receiver)
	origin := helpers.ToByte32(hexutil.MustDecode(task.Origin))
	signature := hexutil.MustDecode(task.Signature)
	amount, err := getAmountOrDefault(transfer.Amount, big.NewInt(1))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid amount: %s", transfer.Amount))
	}

	args := contract.WithdrawArgs{
		Amount:     amount.Uint64(),
		Path:       task.MustParseMerklePath(),
		RecoveryId: signature[64],
		Seeds:      c.solana.BridgeAdminSeed,
		Origin:     origin,
		Signature:  helpers.ToByte64(signature),
	}

	withdrawAddress, _, err := solana.FindProgramAddress([][]byte{origin[:]}, c.solana.BridgeProgramID)
	if err != nil {
		return errors.New("failed to create withdraw address")
	}

	var instruction solana.Instruction
	switch tokenDetails.TokenType {
	case tokenmanager.Type_NATIVE:
		instruction, err = contract.WithdrawNativeInstruction(
			c.solana.BridgeProgramID,
			c.solana.BridgeAdmin,
			solana.PublicKeyFromBytes(receiver),
			withdrawAddress,
			args,
		)
	case tokenmanager.Type_METAPLEX_FT:
		instruction, err = contract.WithdrawFTInstruction(
			c.solana.BridgeProgramID,
			c.solana.BridgeAdmin,
			solana.PublicKeyFromBytes(hexutil.MustDecode(token.TokenAddress)),
			solana.PublicKeyFromBytes(receiver),
			withdrawAddress,
			args,
		)
	case tokenmanager.Type_METAPLEX_NFT:
		instruction, err = contract.WithdrawNFTInstruction(
			c.solana.BridgeProgramID,
			c.solana.BridgeAdmin,
			solana.PublicKeyFromBytes(hexutil.MustDecode(token.TokenId)),
			solana.PublicKeyFromBytes(receiver),
			withdrawAddress,
			args,
		)
	default:
		return errors.Errorf("invalid solana token type: %d", tokenDetails.TokenType)
	}
	if err != nil {
		return errors.Wrap(err, "failed to construct the solana instruction")
	}

	recent, err := c.solana.RPC.GetLatestBlockhash(
		context.Background(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return errors.Wrap(err, "failed to fetch recent blockhash")
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.solana.SubmitterPrivateKey.PublicKey()),
	)
	if err != nil {
		return errors.Wrap(err, "failed to form a solana transaction")
	}

	if _, err = tx.AddSignature(c.solana.SubmitterPrivateKey); err != nil {
		return errors.Wrap(err, "failed to sign a solana transaction")
	}

	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		c.solana.RPC,
		c.solana.WS,
		tx,
	)
	if err != nil {
		return errors.Wrap(err, "failed to submit a solana transaction")
	}

	log.WithField("sig", sig.String()).Debug("successfully submitted transaction")

	return nil
}
