package scheduler

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
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/data/pg"
)

type Service struct {
	client          *http.HTTP
	log             *logan.Entry
	rarimocore      rarimocore.QueryClient
	storage         *pg.Storage
	catchupDisabled bool
}

func NewService(cfg config.Config) *Service {
	return &Service{
		client:          cfg.Tendermint(),
		log:             cfg.Log(),
		rarimocore:      rarimocore.NewQueryClient(cfg.Cosmos()),
		storage:         pg.New(cfg.DB()),
		catchupDisabled: cfg.Relay().CatchupDisabled,
	}
}

func (s *Service) Run(ctx context.Context) {
	if !s.catchupDisabled {
		if err := s.catchup(ctx); err != nil {
			s.log.WithError(err).Error("failed to catchup")
		}
	}

	running.WithBackOff(
		ctx, s.log, "Service", s.run,
		5*time.Second, 5*time.Second, 5*time.Second,
	)
}

func (s *Service) run(ctx context.Context) error {
	s.log.Info("Starting subscription")
	defer s.log.Info("Subscription finished")

	const depositChanSize = 100

	out, err := s.client.Subscribe(
		ctx,
		"Service",
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

func (s *Service) catchup(ctx context.Context) error {
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

func (s *Service) process(
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

func (s *Service) trySave(ctx context.Context, operation rarimocore.Operation) error {
	if operation.OperationType == rarimocore.OpType_IDENTITY_DEFAULT_TRANSFER {
		s.log.WithField("operation_index", operation.Index).Info("Trying to save op")

		op, err := pkg.GetIdentityDefaultTransfer(operation)
		if err != nil {
			return errors.Wrap(err, "failed to parse identity default transfer", logan.F{
				"operation_index": operation.Index,
			})
		}

		err = s.storage.StateQ().UpsertCtx(ctx, &data.State{
			ID:        op.StateHash,
			Operation: operation.Index,
		})

		if err != nil {
			return errors.Wrap(err, "failed to upsert identity default transfer", logan.F{
				"operation_index": operation.Index,
			})
		}
	}

	return nil
}
