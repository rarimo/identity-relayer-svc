package ingester

import (
	"context"
	"github.com/rarimo/identity-relayer-svc/internal/data"
	"github.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/rarimo/identity-relayer-svc/internal/config"
	"github.com/rarimo/identity-relayer-svc/internal/data/pg"
	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type aggregatedIngester struct {
	log        *logan.Entry
	rarimocore rarimocore.QueryClient
	storage    *pg.Storage
}

var _ Processor = &aggregatedIngester{}

func NewAggregatedIngester(cfg config.Config) Processor {
	return &aggregatedIngester{
		log:        cfg.Log(),
		rarimocore: rarimocore.NewQueryClient(cfg.Cosmos()),
		storage:    pg.New(cfg.DB()),
	}
}

func (s *aggregatedIngester) query() string {
	return aggregatedQuery
}

func (s *aggregatedIngester) name() string {
	return "identity-aggregated-ingester"
}

func (s *aggregatedIngester) catchup(ctx context.Context) error {
	s.log.Info("Starting catchup")
	defer s.log.Info("Catchup finished")

	var nextKey []byte

	for {
		operations, err := s.rarimocore.OperationAll(ctx, &rarimocore.QueryAllOperationRequest{Pagination: &query.PageRequest{Key: nextKey}})
		if err != nil {
			panic(err)
		}

		for _, op := range operations.Operation {
			if op.Status == rarimocore.OpStatus_SIGNED && op.OperationType == rarimocore.OpType_IDENTITY_AGGREGATED_TRANSFER {
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

func (s *aggregatedIngester) process(
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

func (s *aggregatedIngester) trySave(ctx context.Context, operation rarimocore.Operation) error {
	if operation.OperationType == rarimocore.OpType_IDENTITY_AGGREGATED_TRANSFER {
		s.log.WithField("operation_index", operation.Index).Info("Trying to save op")

		op, err := pkg.GetIdentityAggregatedTransfer(operation)
		if err != nil {
			return errors.Wrap(err, "failed to parse identity aggregated transfer", logan.F{
				"operation_index": operation.Index,
			})
		}

		err = s.storage.AggregatedQ().UpsertCtx(ctx, &data.Aggregated{
			Gist:      op.GISTHash,
			StateRoot: op.StateRootHash,
			Operation: operation.Index,
		})

		if err != nil {
			return errors.Wrap(err, "failed to upsert identity aggregated transfer", logan.F{
				"operation_index": operation.Index,
			})
		}
	}

	return nil
}
