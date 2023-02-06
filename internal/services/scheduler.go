package services

import (
	"context"
	"time"

	"github.com/adjust/rmq/v5"
	client "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"golang.org/x/exp/slices"

	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/services/relayer"
)

const (
	DepositChanSize = 100
)

type Scheduler interface {
	ScheduleRelays(
		ctx context.Context,
		confirmationID string,
		transferIndexes []string,
	) error
}

type scheduler struct {
	client       *http.HTTP
	log          *logan.Entry
	cosmos       client.ServiceClient
	core         core.Core
	tokenmanager tokenmanager.QueryClient
	rarimocore   rarimocore.QueryClient
	relayQueue   rmq.Queue
}

func NewScheduler(cfg config.Config) Scheduler {
	return newScheduler(cfg)
}

func newScheduler(cfg config.Config) *scheduler {
	return &scheduler{
		client:       cfg.Tendermint(),
		log:          cfg.Log().WithField("service", "scheduler"),
		cosmos:       client.NewServiceClient(cfg.Cosmos()),
		tokenmanager: tokenmanager.NewQueryClient(cfg.Cosmos()),
		rarimocore:   rarimocore.NewQueryClient(cfg.Cosmos()),
		relayQueue:   cfg.Redis().OpenRelayQueue(),
		core:         core.NewCore(cfg),
	}
}

func RunScheduler(cfg config.Config, ctx context.Context) {
	s := newScheduler(cfg)
	running.WithBackOff(
		ctx, s.log, "scheduler", s.run,
		5*time.Second, 5*time.Second, 5*time.Second,
	)
}

func (s *scheduler) run(ctx context.Context) error {
	out, err := s.client.Subscribe(
		ctx,
		"scheduler",
		"tm.event='Tx' AND message.action='create_confirmation'",
		DepositChanSize,
	)
	s.log.Info("listening for confirmations")
	if err != nil {
		return errors.Wrap(err, "can not subscribe")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c := <-out:
			// Delay for indexing tx in core databases
			time.Sleep(time.Second * 5)
			tx, err := s.cosmos.GetTx(ctx, &client.GetTxRequest{Hash: c.Events["tx.hash"][0]})
			if err != nil {
				s.log.WithError(err).Error("error fetching tx by hash")
				continue
			}

			for _, message := range tx.Tx.Body.Messages {
				if message.TypeUrl != "/rarifyprotocol.rarimocore.rarimocore.MsgCreateConfirmation" {
					continue
				}

				msg := rarimocore.MsgCreateConfirmation{}
				if err = msg.Unmarshal(message.Value); err != nil {
					s.log.WithError(err).Error("failed to unmarshal message")
					continue
				}

				if err := s.ScheduleRelays(ctx, msg.Root, msg.Indexes); err != nil {
					s.log.WithError(err).Error("failed to schedule")
				}
			}
		}
	}
}

func (s *scheduler) ScheduleRelays(
	ctx context.Context,
	confirmationID string,
	transferIndexes []string,
) error {
	log := s.log.WithField("merkle_root", confirmationID)
	log.Info("processing a confirmation")

	transfers, err := s.core.GetTransfers(ctx, confirmationID)
	if err != nil {
		return errors.Wrap(err, "failed to get transfers")
	}

	tasks := []data.RelayTask{}
	for _, transfer := range transfers {
		if !slices.Contains(transferIndexes, transfer.Transfer.Origin) {
			continue
		}
		tasks = append(tasks, data.NewRelayTask(transfer, relayer.MaxRetries))
	}

	rawTasks := [][]byte{}
	for _, task := range tasks {
		if slices.Contains(transferIndexes, task.OperationIndex) {
			rawTasks = append(rawTasks, task.Marshal())
		}
	}

	if err := s.relayQueue.PublishBytes(rawTasks...); err != nil {
		return errors.Wrap(err, "failed to publish tasks")
	}

	log.Infof("scheduled %d transfers for relay", len(rawTasks))

	return nil

}
