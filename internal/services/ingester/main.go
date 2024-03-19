package ingester

import (
	"context"
	"fmt"
	"time"

	"github.com/rarimo/identity-relayer-svc/internal/config"
	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	stateQuery = "tm.event='Tx' AND operation_signed.operation_type='IDENTITY_STATE_TRANSFER'"
	gistQuery  = "tm.event='Tx' AND operation_signed.operation_type='IDENTITY_GIST_TRANSFER'"

	// Iden3 Issuer
	aggregatedQuery = "tm.event='Tx' AND operation_signed.operation_type='IDENTITY_AGGREGATED_TRANSFER'"
)

type Service struct {
	Processor
	log             *logan.Entry
	client          *http.HTTP
	catchupDisabled bool
}

type Processor interface {
	catchup(ctx context.Context) error
	process(ctx context.Context, confirmationID string) error
	query() string
	name() string
}

func NewService(cfg config.Config, processor Processor) *Service {
	return &Service{
		Processor:       processor,
		log:             cfg.Log(),
		client:          cfg.Tendermint(),
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
		ctx, s.log, s.Processor.name(), s.run,
		5*time.Second, 5*time.Second, 5*time.Second,
	)
}

func (s *Service) run(ctx context.Context) error {
	s.log.Info("Starting subscription")
	defer s.log.Info("Subscription finished")

	const depositChanSize = 100

	out, err := s.client.Subscribe(
		ctx,
		s.Processor.name(),
		s.Processor.query(),
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
			s.log.Infof("New confirmation found %s", confirmation)

			if err := s.process(ctx, confirmation); err != nil {
				s.log.WithError(err).Error("failed to process confirmation")
			}
		}
	}
}
