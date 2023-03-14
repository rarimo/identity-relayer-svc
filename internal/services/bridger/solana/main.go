package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-redis/redis/v8"
	"github.com/olegfomenko/solana-go"
	"github.com/olegfomenko/solana-go/rpc"
	confirm "github.com/olegfomenko/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger/bridge"
	"gitlab.com/rarimo/relayer-svc/internal/utils"
	"gitlab.com/rarimo/solana-program-go/contract"
)

const solanaEstimateTTLMinutes = 30

type solanaBridger struct {
	log          *logan.Entry
	rarimocore   rarimocore.QueryClient
	tokenmanager tokenmanager.QueryClient
	solana       *config.Solana
	tollbooth    *config.Tollbooth
	redis        *redis.Client
}

func NewSolanaBridger(cfg config.Config) bridge.Bridger {
	return &solanaBridger{
		log:          cfg.Log().WithField("service", "solana_bridger"),
		rarimocore:   rarimocore.NewQueryClient(cfg.Cosmos()),
		tokenmanager: tokenmanager.NewQueryClient(cfg.Cosmos()),
		solana:       cfg.Solana(),
		redis:        cfg.Redis().Client(),
		tollbooth:    cfg.Tollbooth(),
	}
}

func (b *solanaBridger) Withdraw(
	ctx context.Context,
	transfer core.TransferDetails,
) error {
	log := b.log.WithField("op_id", transfer.Origin)
	withdrawn, err := b.isAlreadyWithdrawn(ctx, transfer)
	if err != nil {
		return errors.Wrap(err, "failed to check if the transfer is withdrawn")
	}
	if withdrawn {
		return bridge.ErrAlreadyWithdrawn
	}

	tx, err := b.makeWitdrawTx(ctx, transfer)
	if err != nil {
		return errors.Wrap(err, "failed to call the withdraw method")
	}
	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		b.solana.RPC,
		b.solana.WS,
		tx,
	)
	if err != nil {
		return errors.Wrap(err, "failed to submit a solana transaction")
	}

	log.WithFields(logan.F{"sig": sig.String()}).Info("successfully submitted transaction")

	return nil
}

func (b *solanaBridger) EstimateRelayFee(
	ctx context.Context,
	transfer core.TransferDetails,
) (bridge.FeeEstimate, error) {
	withdrawn, err := b.isAlreadyWithdrawn(ctx, transfer)
	if err != nil {
		return bridge.FeeEstimate{}, errors.Wrap(err, "failed to check if the transfer is withdrawn")
	}
	if withdrawn {
		return bridge.FeeEstimate{}, bridge.ErrAlreadyWithdrawn
	}

	estimate, err := b.getSavedEstimate(ctx, transfer.Origin)
	if err != nil {
		return bridge.FeeEstimate{}, errors.Wrap(err, "failed to get saved estimate")
	}
	if estimate != nil {
		return *estimate, nil
	}

	tx, err := b.makeWitdrawTx(ctx, transfer)
	if err != nil {
		return bridge.FeeEstimate{}, errors.Wrap(err, "failed to create the withdraw method")
	}

	fee, err := b.solana.RPC.GetFeeForMessage(ctx, tx.Message.ToBase64(), "")
	if err != nil {
		return bridge.FeeEstimate{}, errors.Wrap(err, "failed to get the fee for the message")
	}
	if fee.Value == nil {
		return bridge.FeeEstimate{}, errors.New("fee is nil")
	}
	gasEstimate := big.NewInt(0).SetUint64(*fee.Value)

	createdAt := time.Now()
	feeEstimate := bridge.FeeEstimate{
		TransferID:  transfer.Transfer.Origin,
		GasEstimate: gasEstimate,
		GasToken:    b.tollbooth.Solana.GasTokenTicker,
		// NOTE: we use a hardcoded value for the relay fee because it's too miniscule to bother
		FeeAmount:       b.tollbooth.Solana.RelayFee,
		FeeToken:        b.tollbooth.FeeTokenTicker,
		FeeTokenAddress: b.tollbooth.Solana.FeeTokenMint.String(),
		FromChain:       transfer.Transfer.FromChain,
		ToChain:         transfer.Transfer.ToChain,
		CreatedAt:       createdAt,
		ExpiresAt:       createdAt.Add(solanaEstimateTTLMinutes * time.Minute),
	}

	if err := b.saveEstimate(ctx, &feeEstimate); err != nil {
		return bridge.FeeEstimate{}, errors.Wrap(err, "failed to save the estimate")
	}

	return feeEstimate, nil
}

func (b *solanaBridger) makeWitdrawTx(
	ctx context.Context,
	transfer core.TransferDetails,
) (*solana.Transaction, error) {
	receiver := hexutil.MustDecode(transfer.Transfer.Receiver)
	origin := utils.ToByte32(hexutil.MustDecode(transfer.Origin))
	signature := hexutil.MustDecode(transfer.Signature)
	amount, err := utils.GetAmountOrDefault(transfer.Transfer.Amount, big.NewInt(1))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("invalid amount: %s", transfer.Transfer.Amount))
	}

	args := contract.WithdrawArgs{
		Amount:     amount.Uint64(),
		Path:       transfer.MerklePath,
		RecoveryId: signature[64],
		Seeds:      b.solana.BridgeAdminSeed,
		Origin:     origin,
		Signature:  utils.ToByte64(signature),
	}

	withdrawAddress, _, err := solana.FindProgramAddress([][]byte{origin[:]}, b.solana.BridgeProgramID)
	if err != nil {
		return nil, errors.New("failed to create withdraw address")
	}

	var instruction solana.Instruction
	switch transfer.TokenDetails.TokenType {
	case tokenmanager.Type_NATIVE:
		instruction, err = contract.WithdrawNativeInstruction(
			b.solana.BridgeProgramID,
			b.solana.BridgeAdmin,
			solana.PublicKeyFromBytes(receiver),
			withdrawAddress,
			args,
		)
	case tokenmanager.Type_METAPLEX_FT:
		tokenAddress := hexutil.MustDecode(transfer.TokenDetails.TokenAddress)
		instruction, err = contract.WithdrawFTInstruction(
			b.solana.BridgeProgramID,
			b.solana.BridgeAdmin,
			solana.PublicKeyFromBytes(tokenAddress),
			solana.PublicKeyFromBytes(receiver),
			withdrawAddress,
			args,
		)
	case tokenmanager.Type_METAPLEX_NFT:
		tokenID := hexutil.MustDecode(transfer.TokenDetails.TokenId)
		instruction, err = contract.WithdrawNFTInstruction(
			b.solana.BridgeProgramID,
			b.solana.BridgeAdmin,
			solana.PublicKeyFromBytes(tokenID),
			solana.PublicKeyFromBytes(receiver),
			withdrawAddress,
			args,
		)
	default:
		return nil, errors.Errorf("invalid solana token type: %d", transfer.TokenDetails.TokenType)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct the solana instruction")
	}

	recent, err := b.solana.RPC.GetLatestBlockhash(
		context.Background(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch recent blockhash")
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(b.solana.SubmitterPrivateKey.PublicKey()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to form a solana transaction")
	}

	if _, err = tx.AddSignature(b.solana.SubmitterPrivateKey); err != nil {
		return nil, errors.Wrap(err, "failed to sign a solana transaction")
	}

	return tx, nil
}

func (b *solanaBridger) getSavedEstimate(ctx context.Context, transferID string) (*bridge.FeeEstimate, error) {
	resp := b.redis.Get(ctx, transferID)
	if resp.Err() != nil {
		if resp.Err() == redis.Nil {
			return nil, nil
		}
		return nil, errors.Wrap(resp.Err(), "failed to get the estimate")
	}

	raw, err := resp.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read estimate from redis")
	}

	var estimate bridge.FeeEstimate
	if err := json.Unmarshal(raw, &estimate); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal estimate")
	}

	return &estimate, nil
}

func (b *solanaBridger) saveEstimate(ctx context.Context, estimate *bridge.FeeEstimate) error {
	raw, err := json.Marshal(estimate)
	if err != nil {
		return errors.Wrap(err, "failed to marshal estimate")
	}

	if resp := b.redis.Set(ctx, estimate.TransferID, raw, solanaEstimateTTLMinutes*time.Minute); resp.Err() != nil {
		return errors.Wrap(resp.Err(), "failed to set the cursor")
	}

	return nil
}

func (b solanaBridger) isAlreadyWithdrawn(ctx context.Context, transfer core.TransferDetails) (bool, error) {
	origin := utils.ToByte32(hexutil.MustDecode(transfer.Origin))
	withdrawAddress, _, err := solana.FindProgramAddress([][]byte{origin[:]}, b.solana.BridgeProgramID)
	if err != nil {
		return false, errors.New("failed to create withdraw address")
	}
	_, err = b.solana.RPC.GetAccountInfoWithOpts(
		ctx, withdrawAddress,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentType(rpc.ConfirmationStatusProcessed),
		},
	)
	if errors.Cause(err) == rpc.ErrNotFound {
		// has not been withdrawn yet
	} else if err != nil {
		return false, errors.Wrap(err, "failed to get withdraw account")
	} else {
		return true, nil
	}

	return false, nil
}
