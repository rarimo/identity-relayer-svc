package ingester

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/rarimo/identity-relayer-svc/internal/config"
	"github.com/rarimo/identity-relayer-svc/internal/data"
	"github.com/rarimo/identity-relayer-svc/internal/data/pg"
	"github.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type stateIngester struct {
	log        *logan.Entry
	rarimocore rarimocore.QueryClient
	storage    *pg.Storage
	filter     func(string) bool
}

var _ Processor = &stateIngester{}

func NewStateIngester(cfg config.Config) Processor {
	filtrationDisabled := cfg.Relay().DisableFiltration
	allowList := mp(cfg.Relay().IssuerID)

	filter := func(id string) bool {
		if filtrationDisabled {
			return true
		}
		_, ok := allowList[id]
		return ok
	}

	return &stateIngester{
		log:        cfg.Log(),
		rarimocore: rarimocore.NewQueryClient(cfg.Cosmos()),
		storage:    pg.New(cfg.DB()),
		filter:     filter,
	}
}

func (s *stateIngester) query() string {
	return stateQuery
}

func (s *stateIngester) name() string {
	return "identity-state-ingester"
}

func (s *stateIngester) catchup(ctx context.Context) error {
	s.log.Info("Starting catchup")
	defer s.log.Info("Catchup finished")

	var nextKey []byte

	for {
		operations, err := s.rarimocore.OperationAll(ctx, &rarimocore.QueryAllOperationRequest{Pagination: &query.PageRequest{Key: nextKey}})
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

func (s *stateIngester) process(
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

func (s *stateIngester) trySave(ctx context.Context, operation rarimocore.Operation) error {
	if operation.OperationType == rarimocore.OpType_IDENTITY_STATE_TRANSFER {
		s.log.WithField("operation_index", operation.Index).Info("Trying to save op")

		op, err := pkg.GetIdentityStateTransfer(operation)
		if err != nil {
			return errors.Wrap(err, "failed to parse identity default transfer", logan.F{
				"operation_index": operation.Index,
			})
		}

		if !s.filter(op.Id) {
			s.log.WithField("operation_index", operation.Index).Info("Issuer ID is not supported")
			return nil
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
