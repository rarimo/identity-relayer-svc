package solana

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-redis/redis/v8"
	"github.com/olegfomenko/solana-go"
	"github.com/olegfomenko/solana-go/rpc"
	"github.com/olegfomenko/solana-go/rpc/ws"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"

	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/services"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger/bridge"
	sb "gitlab.com/rarimo/relayer-svc/internal/services/bridger/solana"
)

const solanaListenerName = "solana_listener"

func RunSolanaListener(cfg config.Config, ctx context.Context) {
	chain := cfg.Solana()
	log := cfg.Log().WithField("runner", solanaListenerName)
	log.Info("starting listener")

	listener := solanaListener{
		log:             log,
		redis:           cfg.Redis().Client(),
		chain:           chain,
		core:            core.NewCore(cfg),
		cursorKey:       fmt.Sprintf("%s:tollbooth-cursor", "solana"),
		scheduler:       services.NewScheduler(cfg),
		bridge:          sb.NewSolanaBridger(cfg),
		bridgerProvider: bridger.NewBridgerProvider(cfg),
		tollbooth:       cfg.Tollbooth(),
	}

	running.UntilSuccess(ctx, log, solanaListenerName, func(ctx context.Context) (bool, error) {
		if err := listener.catchup(ctx); err != nil {
			return false, errors.Wrap(err, "failed to catchup")
		}

		return true, nil
	}, 1*time.Second, 1*time.Second)

	log.Info("starting subscription")
	running.WithBackOff(ctx, log, solanaListenerName,
		listener.subscribe,
		5*time.Second, 5*time.Second, 5*time.Second,
	)

	log.Info("finished listener")
}

type solanaListener struct {
	log             *logan.Entry
	redis           *redis.Client
	chain           *config.Solana
	core            core.Core
	scheduler       services.Scheduler
	cursorKey       string
	bridge          bridge.Bridger
	bridgerProvider bridger.BridgerProvider
	tollbooth       *config.Tollbooth
}

func (l *solanaListener) catchupFrom(
	ctx context.Context,
	cursor solana.Signature,
	targetCursor solana.Signature,
) (solana.Signature, error) {
	l.log.WithField("cursor", cursor).Info("Catching up history")

	signatures, err := l.chain.RPC.GetSignaturesForAddressWithOpts(
		ctx,
		l.tollbooth.Solana.TollboothAccount,
		&rpc.GetSignaturesForAddressOpts{
			Before:     cursor,
			Commitment: rpc.CommitmentFinalized,
		},
	)

	if err != nil {
		return solana.Signature{}, errors.Wrap(err, "error getting txs")
	}

	for _, sig := range signatures {
		l.log.Debug("Checking tx: " + sig.Signature.String())
		if err := l.handleTransaction(ctx, sig.Signature); err != nil {
			l.log.
				WithError(err).
				WithField("signature", sig.Signature).
				Error("failed to process a transaction")
		}

		if targetCursor.Equals(sig.Signature) {
			return sig.Signature, nil
		}
	}

	if len(signatures) == 0 {
		return cursor, nil
	}

	return signatures[len(signatures)-1].Signature, nil
}

func (l *solanaListener) catchup(ctx context.Context) (err error) {
	l.log.Info("starting catchup")
	cursorFrom, err := l.getCursor(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting cursor")
	}
	if cursorFrom.Equals(solana.Signature{}) {
		l.log.Info("no cursor found, starting from genesis")
	}

	// Catchup goes backwards from the latest signature in the chain to the saved cursor.
	// We need to get the latest signature first to save it as the cursor after the catchup.
	one := 1
	signatures, err := l.chain.RPC.GetSignaturesForAddressWithOpts(
		ctx,
		l.tollbooth.Solana.FeeTokenMint,
		&rpc.GetSignaturesForAddressOpts{
			Commitment: rpc.CommitmentFinalized,
			Limit:      &one,
		},
	)
	if err != nil {
		return errors.Wrap(err, "error getting the last signature")
	}
	if len(signatures) == 0 {
		l.log.Info("nothing to catchup")
		return nil
	}
	cursorTo := signatures[0].Signature

	var start solana.Signature
	for {
		last, err := l.catchupFrom(ctx, start, cursorFrom)
		if err != nil {
			return errors.Wrap(err, "failed to catch up history", logan.F{
				"cursor": cursorFrom,
			})
		}

		if cursorFrom.Equals(last) || last.Equals(start) {
			break
		}

		start = last
	}

	if err := l.setCursor(ctx, cursorTo); err != nil {
		return errors.Wrap(err, "error setting cursor")
	}
	l.log.WithField("cursor", cursorTo).Info("Catchup finished")

	return nil
}

func (l *solanaListener) subscribe(ctx context.Context) error {
	sub, err := l.chain.WS.LogsSubscribeMentions(
		l.tollbooth.Solana.TollboothAccount,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return errors.Wrap(err, "error subscribing to the program logs")
	}

	defer sub.Unsubscribe()

	runnerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	rec := make(chan *ws.LogResult, 1)
	go func() {
		for runnerCtx.Err() == nil {
			logEntry, err := sub.Recv()
			if err != nil {
				panic(errors.Wrap(err, "failed to receive the message"))
			}
			rec <- logEntry
		}
	}()

	for {
		select {
		case <-runnerCtx.Done():
			return runnerCtx.Err()
		case logEntry := <-rec:
			if err := l.handleTransaction(ctx, logEntry.Value.Signature); err != nil {
				l.log.
					WithError(err).
					WithField("signature", logEntry.Value.Signature).
					Error("failed to process a transaction")
			}

			if err := l.setCursor(ctx, logEntry.Value.Signature); err != nil {
				l.log.
					WithError(err).
					WithField("signature", logEntry.Value.Signature).
					Error("failed to set cursor")
			}
		}
	}
}

func (l *solanaListener) handleTransaction(ctx context.Context, signature solana.Signature) error {
	tx, err := l.chain.RPC.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{
		Encoding: solana.EncodingBase64,
	})
	if err != nil {
		return errors.Wrap(err, "error getting transaction from solana")
	}
	if tx == nil {
		return errors.New("transaction not found")
	}
	event, err := FindFeePaidEvent(tx.Meta.LogMessages)
	if err != nil {
		return errors.Wrap(err, "error finding fee paid event")
	}

	confirmationID := hexutil.Encode(event.ConfirmationId[:])
	transferID := hexutil.Encode(event.TransferId[:])

	transfer, err := l.core.GetTransfer(ctx, confirmationID, transferID)
	if err != nil {
		return errors.Wrap(err, "failed to get transfer")
	}

	targetBridger := l.bridgerProvider.GetBridger(transfer.Transfer.To.Chain)
	estimate, err := targetBridger.EstimateRelayFee(ctx, *transfer)
	if errors.Cause(err) == bridge.ErrAlreadyWithdrawn {
		l.log.Info("transfer already withdrawn, skipping")
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to estimate relay fee")
	}

	if event.FeeTokenMint.String() != estimate.FeeTokenAddress {
		return errors.New("fee token is not the same as the estimated fee token")
	}
	if big.NewInt(0).SetUint64(event.Amount).Cmp(estimate.FeeAmount) < 0 {
		return errors.New("fee paid is less than the estimated fee")
	}

	if err := l.scheduler.ScheduleRelays(ctx, confirmationID, []string{transferID}); err != nil {
		panic(errors.Wrap(err, "failed to schedule the transfers for relay"))
	}

	return nil
}

func (l *solanaListener) getCursor(ctx context.Context) (solana.Signature, error) {
	resp := l.redis.Get(ctx, l.cursorKey)
	if resp.Err() != nil {
		if resp.Err() == redis.Nil {
			return solana.Signature{}, nil
		}
		return solana.Signature{}, errors.Wrap(resp.Err(), "failed to get cursor")
	}

	rawCursor := resp.Val()
	if rawCursor == "" {
		return solana.Signature{}, nil
	}

	cursor, err := solana.SignatureFromBase58(rawCursor)
	if err != nil {
		l.log.WithError(err).WithField("raw", rawCursor).Error("failed to parse the cursor value")
		return solana.Signature{}, nil
	}

	return cursor, nil
}

func (l *solanaListener) setCursor(ctx context.Context, cursor solana.Signature) error {
	if resp := l.redis.Set(ctx, l.cursorKey, cursor.String(), 0); resp.Err() != nil {
		return errors.Wrap(resp.Err(), "failed to set cursor")
	}

	return nil
}
