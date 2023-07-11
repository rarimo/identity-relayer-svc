package services

import (
	"context"
	"fmt"
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
		operationIndexes []string,
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
		"tm.event='Tx' AND operation_signed.operation_type='IDENTITY_DEFAULT_TRANSFER'",
		DepositChanSize,
	)
	s.log.Info("Listening for signed identities")

	if err != nil {
		return errors.Wrap(err, "can not subscribe")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c := <-out:
			s.log.Info("New confirmation found")

			index := c.Events[fmt.Sprintf("%s.%s", rarimocore.EventTypeOperationSigned, rarimocore.AttributeKeyOperationId)][0]
			confirmation := c.Events[fmt.Sprintf("%s.%s", rarimocore.EventTypeOperationSigned, rarimocore.AttributeKeyConfirmationId)][0]

			s.log.Infof("New operation found index=%s, conf=%s", index, confirmation)

			if err := s.ScheduleRelays(ctx, confirmation, []string{index}); err != nil {
				s.log.WithError(err).Error("failed to schedule")
			}
		}
	}
}

func (s *scheduler) ScheduleRelays(
	ctx context.Context,
	confirmationID string,
	operationIndexes []string,
) error {
	log := s.log.WithField("merkle_root", confirmationID)
	log.Info("processing a confirmation")

	rawTasks := [][]byte{}

	for _, index := range operationIndexes {
		operation, err := s.rarimocore.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: index})
		if err != nil {
			return errors.Wrap(err, "failed to get operation")
		}

		switch operation.Operation.OperationType {
		case rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER:
			transfer, err := s.core.GetIdentityDefaultTransfer(ctx, confirmationID, operation.Operation.Index)
			if err != nil {
				return errors.Wrap(err, "failed to get identity default transfer", logan.F{
					"confirmation_id": confirmationID,
					"operation_index": operation.Operation.Index,
				})
			}

			if !slices.Contains(operationIndexes, transfer.OpIndex) {
				continue
			}

			rawTasks = append(rawTasks, data.NewRelayIdentityTransferTask(*transfer, relayer.MaxRetries).Marshal())
		}
	}

	if err := s.relayQueue.PublishBytes(rawTasks...); err != nil {
		return errors.Wrap(err, "failed to publish tasks")
	}

	log.Infof("scheduled %d transfers for relay", len(rawTasks))

	return nil

}
