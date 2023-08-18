package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/data/pg"
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
	tokenmanager tokenmanager.QueryClient
	rarimocore   rarimocore.QueryClient
	storage      *pg.Storage
}

func newScheduler(cfg config.Config) *scheduler {
	return &scheduler{
		client:       cfg.Tendermint(),
		log:          cfg.Log().WithField("service", "scheduler"),
		tokenmanager: tokenmanager.NewQueryClient(cfg.Cosmos()),
		rarimocore:   rarimocore.NewQueryClient(cfg.Cosmos()),
		storage:      pg.New(cfg.DB()),
	}
}

func RunScheduler(cfg config.Config, ctx context.Context) {
	s := newScheduler(cfg)

	if !cfg.Relay().CatchupDisabled {
		if err := s.catchup(ctx); err != nil {

		}
	}

	running.WithBackOff(
		ctx, s.log, "scheduler", s.run,
		5*time.Second, 5*time.Second, 5*time.Second,
	)
}

func (s *scheduler) run(ctx context.Context) error {
	s.log.Info("Starting subscription")
	defer s.log.Info("Subscription finished")

	const depositChanSize = 100

	out, err := s.client.Subscribe(
		ctx,
		"scheduler",
		"tm.event='Tx' AND operation_signed.operation_type='IDENTITY_DEFAULT_TRANSFER'",
		depositChanSize,
	)

	if err != nil {
		return errors.Wrap(err, "can not subscribe")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c := <-out:
			s.log.Info("New confirmation found")
			confirmation := c.Events[fmt.Sprintf("%s.%s", rarimocore.EventTypeOperationSigned, rarimocore.AttributeKeyConfirmationId)][0]
			s.log.Infof("New confirmation for identity found %s", confirmation)

			if err := s.process(ctx, confirmation); err != nil {
				s.log.WithError(err).Error("failed to process confirmation")
			}
		}
	}
}

func (s *scheduler) catchup(ctx context.Context) error {
	s.log.Info("Starting catchup")
	defer s.log.Info("Catchup finished")

	var nextKey []byte

	for {
		operations, err := s.rarimocore.OperationAll(context.TODO(), &rarimocore.QueryAllOperationRequest{Pagination: &query.PageRequest{Key: nextKey}})
		if err != nil {
			panic(err)
		}

		for _, op := range operations.Operation {
			if op.Status == rarimocore.OpStatus_SIGNED && op.OperationType == rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER {
				if err := s.trySave(ctx, op); err != nil {
					return err
				}
			}
		}

		nextKey = operations.Pagination.NextKey
		if nextKey == nil {
			return nil
		}
	}
}

func (s *scheduler) process(
	ctx context.Context,
	confirmationID string,
) error {
	log := s.log.WithField("confirmation_id", confirmationID)
	log.Info("processing a confirmation")

	rawConf, err := s.rarimocore.Confirmation(ctx, &rarimocore.QueryGetConfirmationRequest{Root: confirmationID})
	if err != nil {
		return errors.Wrap(err, "failed to get confirmation", logan.F{
			"confirmation_id": confirmationID,
		})
	}

	for _, index := range rawConf.Confirmation.Indexes {
		operation, err := s.rarimocore.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: index})
		if err != nil {
			return errors.Wrap(err, "failed to get operation", logan.F{
				"operation_index": operation.Operation.Index,
			})
		}

		if err := s.trySave(ctx, operation.Operation); err != nil {
			return err
		}
	}

	return nil

}

func (s *scheduler) trySave(ctx context.Context, operation rarimocore.Operation) error {
	if operation.OperationType == rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER {
		op, err := pkg.GetIdentityDefaultTransfer(operation)
		if err != nil {
			return errors.Wrap(err, "failed to parse identity default transfer", logan.F{
				"operation_index": operation.Index,
			})
		}

		err = s.storage.StateQ().UpsertCtx(ctx, &data.State{
			ID:        op.StateHash,
			Operation: op.Id,
		})

		if err != nil {
			return errors.Wrap(err, "failed to upsert identity default transfer", logan.F{
				"operation_index": operation.Index,
			})
		}
	}

	return nil
}
