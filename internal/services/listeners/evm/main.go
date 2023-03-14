package evm

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-redis/redis/v8"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"

	gobind "gitlab.com/rarimo/evm-tollbooth/gobind/contracts"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/services"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger/bridge"
)

func RunEVMListener(ctx context.Context, cfg config.Config, chainName string) {
	chain, ok := cfg.EVM().GetChainByName(chainName)
	if !ok {
		panic(fmt.Errorf("unknown EVM chain: %s", chainName))
	}

	runnerName := fmt.Sprintf("%s_listener", chainName)
	log := cfg.Log().WithField("runner", runnerName)
	log.Info("starting listener")

	handler, err := gobind.NewTollbooth(cfg.Tollbooth().GetEVMConfig(chainName).TollboothAddress, chain.RPC)
	if err != nil {
		panic(errors.Wrap(err, "failed to init native handler"))
	}

	bridgerProvider := bridger.NewBridgerProvider(cfg)
	listener := evmListener{
		log:             log,
		redis:           cfg.Redis().Client(),
		chain:           chain,
		handler:         handler,
		core:            core.NewCore(cfg),
		cursorKey:       fmt.Sprintf("%s:tollbooth-cursor", chainName),
		scheduler:       services.NewScheduler(cfg),
		bridge:          bridgerProvider.GetBridger(chainName),
		bridgerProvider: bridgerProvider,
	}

	cursor, err := listener.getCursor(ctx)
	if err != nil {
		panic(errors.Wrap(err, "failed to get the cursor"))
	}

	log.Info("starting catchup", logan.F{
		"cursor": fmt.Sprintf("%+v", cursor),
	})

	var currentCursor uint64
	running.UntilSuccess(ctx, log, runnerName, func(ctx context.Context) (bool, error) {
		cursor, err := listener.catchup(ctx, cursor)
		if err != nil {
			return false, errors.Wrap(err, "failed to catchup")
		}

		currentCursor = cursor
		return true, nil
	}, 1*time.Second, 1*time.Second)

	log.Info("catchup finished", logan.F{
		"cursor": fmt.Sprintf("%+v", currentCursor),
	})

	log.Info("starting subscription")
	running.WithBackOff(ctx, log, runnerName,
		listener.subscription(currentCursor),
		5*time.Second, 5*time.Second, 5*time.Second)

	log.Info("finished listener")
}

type evmListener struct {
	log             *logan.Entry
	redis           *redis.Client
	handler         *gobind.Tollbooth
	chain           *config.EVMChain
	core            core.Core
	scheduler       services.Scheduler
	cursorKey       string
	bridge          bridge.Bridger
	bridgerProvider bridger.BridgerProvider
}

func (l *evmListener) catchup(ctx context.Context, cursor uint64) (nextCursor uint64, err error) {
	iter, err := l.handler.FilterFeePaid(&bind.FilterOpts{
		Start:   cursor,
		Context: ctx,
	})

	if err != nil {
		return 0, errors.Wrap(err, "failed to init native iterator")
	}

	for iter.Next() {
		e := iter.Event

		if e == nil {
			l.log.Error("got nil event")
			continue
		}

		l.log.Debug("got event", logan.F{
			"raw": fmt.Sprintf("%+v", e),
		})

		cursor, err = l.handleFeePaid(ctx, e, cursor)
		if err != nil {
			l.log.WithError(err).Error("failed to process event")
			continue
		}
	}

	return cursor, nil
}

func (l *evmListener) subscription(cursor uint64) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		sink := make(chan *gobind.TollboothFeePaid, 10)
		defer close(sink)

		sub, err := l.handler.WatchFeePaid(
			&bind.WatchOpts{
				Context: ctx,
			},
			sink,
		)
		if err != nil {
			panic(errors.Wrap(err, "failed to init subscription"))
		}

		defer sub.Unsubscribe()

	subscriptionPoll:
		for {
			select {
			case <-ctx.Done():
				break subscriptionPoll
			case err := <-sub.Err():
				l.log.WithError(err).Error("subscription error")
				break subscriptionPoll
			case e := <-sink:
				if e == nil {
					l.log.Error("got nil event")
					continue
				}

				l.log.Debug("got event", logan.F{
					"raw": e,
				})

				cursor, err = l.handleFeePaid(ctx, e, cursor)
				if err != nil {
					l.log.WithError(err).Error("failed to process event")
				}
			}
		}

		return nil
	}
}

func (l *evmListener) handleFeePaid(ctx context.Context, e *gobind.TollboothFeePaid, cursor uint64) (uint64, error) {
	confirmationID := hexutil.Encode(e.ConfirmationId.Bytes())
	transferID := hexutil.Encode(e.TransferId.Bytes())

	transfer, err := l.core.GetTransfer(ctx, confirmationID, transferID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get transfer")
	}
	targetBridger := l.bridgerProvider.GetBridger(transfer.Transfer.ToChain)
	estimate, err := targetBridger.EstimateRelayFee(ctx, *transfer)
	if errors.Cause(err) == bridge.ErrAlreadyWithdrawn {
		l.log.Info("transfer already withdrawn, skipping")
		return e.Raw.BlockNumber, nil
	}
	if err != nil {
		return 0, errors.Wrap(err, "failed to estimate relay fee")
	}
	if e.FeeToken.Hex() != estimate.FeeToken {
		return 0, errors.New("fee token is not the same as the estimated fee token")
	}
	if e.TotalAmount.Cmp(estimate.FeeAmount) < 0 {
		return 0, errors.New("fee paid is less than the estimated fee")
	}

	if err := l.scheduler.ScheduleRelays(ctx, confirmationID, []string{transferID}); err != nil {
		panic(errors.Wrap(err, "failed to schedule the transfers for relay"))
	}

	if resp := l.redis.Set(ctx, l.cursorKey, e.Raw.BlockNumber, 0); resp.Err() != nil {
		panic(errors.Wrap(resp.Err(), "failed to set the cursor"))
	}

	return cursor, nil
}

func (l *evmListener) getCursor(ctx context.Context) (uint64, error) {
	resp := l.redis.Get(ctx, l.cursorKey)
	if resp.Err() == redis.Nil {
		return 0, nil
	} else if resp.Err() != nil {
		return 0, errors.Wrap(resp.Err(), "failed to get cursor")
	}

	cursor, err := resp.Uint64()
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse the cursor value", logan.F{
			"raw": resp.String(),
		})
	}

	return cursor, nil
}
