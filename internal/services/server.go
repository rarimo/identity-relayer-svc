package services

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarimo/relayer-svc/docs"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data/pg"
	"gitlab.com/rarimo/relayer-svc/internal/services/relayer"
	"gitlab.com/rarimo/relayer-svc/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerImpl struct {
	types.UnimplementedServiceServer
	log      *logan.Entry
	listener net.Listener
	relayer  *relayer.Service
	storage  *pg.Storage
}

func NewServer(cfg config.Config) *ServerImpl {
	return &ServerImpl{
		log:      cfg.Log(),
		listener: cfg.Listener(),
		relayer:  relayer.NewService(cfg),
		storage:  pg.New(cfg.DB()),
	}
}

func (s *ServerImpl) Run() error {
	grpcGatewayRouter := runtime.NewServeMux()
	httpRouter := http.NewServeMux()

	err := types.RegisterServiceHandlerServer(context.Background(), grpcGatewayRouter, s)
	if err != nil {
		panic(err)
	}

	docs.RegisterOpenAPIService("relayer-svc", httpRouter)
	httpRouter.Handle("/", grpcGatewayRouter)
	return http.Serve(s.listener, httpRouter)
}

var _ types.ServiceServer = &ServerImpl{}

func (s *ServerImpl) Relay(ctx context.Context, req *types.MsgRelayRequest) (*types.MsgRelayResponse, error) {
	tx, err := s.relayer.Relay(ctx, req.State, req.Chain)

	if err != nil {
		switch errors.Cause(err) {
		case relayer.ErrEntryNotFound, relayer.ErrChainNotFound:
			return nil, status.Errorf(codes.NotFound, err.Error())
		case relayer.ErrAlreadySubmitted:
			return nil, status.Errorf(
				codes.InvalidArgument,
				"can not resubmit state transition tx for state: %s on chain: %s", req.State, req.Chain,
			)
		default:
			s.log.WithError(err).Error("got internal error while processing relay request")
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}

	return &types.MsgRelayResponse{Tx: tx}, nil
}

func (s *ServerImpl) Relays(ctx context.Context, req *types.MsgRelaysRequest) (*types.MsgRelaysResponse, error) {
	relays, err := s.relayer.Relays(ctx, req.State)
	if err != nil {
		switch errors.Cause(err) {
		case relayer.ErrEntryNotFound:
			return nil, status.Errorf(codes.NotFound, err.Error())
		default:
			s.log.WithError(err).Error("got internal error while processing relay request")
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}

	resp := &types.MsgRelaysResponse{Relays: make([]*types.Transition, 0, len(relays))}

	for _, r := range relays {
		resp.Relays = append(resp.Relays, &types.Transition{
			Chain: r.Chain,
			Tx:    r.Tx,
		})
	}

	return resp, nil
}
