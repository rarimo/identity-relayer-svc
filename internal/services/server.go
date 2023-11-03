package services

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rarimo/identity-relayer-svc/docs"
	"github.com/rarimo/identity-relayer-svc/internal/config"
	"github.com/rarimo/identity-relayer-svc/internal/data/pg"
	"github.com/rarimo/identity-relayer-svc/internal/services/relayer"
	"github.com/rarimo/identity-relayer-svc/internal/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

func (s *ServerImpl) StateRelay(ctx context.Context, req *types.MsgStateRelayRequest) (*types.MsgRelayResponse, error) {
	if req.Body == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty body")
	}

	s.log.Debugf("Relay request: State - %s; chain - %s", req.Body.Hash, req.Body.Chain)

	tx, err := s.relayer.StateRelay(ctx, req.Body.Hash, req.Body.Chain)

	if err != nil {
		s.log.WithError(err).Debugf("Request failed")

		switch errors.Cause(err) {
		case relayer.ErrEntryNotFound, relayer.ErrChainNotFound:
			return nil, status.Errorf(codes.NotFound, err.Error())
		case relayer.ErrAlreadySubmitted:
			return nil, status.Errorf(
				codes.InvalidArgument,
				"can not resubmit state transition tx for state: %s on chain: %s", req.Body.Hash, req.Body.Chain,
			)
		default:
			s.log.WithError(err).Error("Got internal error while processing relay request")
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}

	return &types.MsgRelayResponse{Tx: tx}, nil
}

func (s *ServerImpl) GISTRelay(ctx context.Context, req *types.MsgGISTRelayRequest) (*types.MsgRelayResponse, error) {
	if req.Body == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty body")
	}

	s.log.Infof("Relay request: GIST - %s; chain - %s", req.Body.Hash, req.Body.Chain)

	tx, err := s.relayer.GistRelay(ctx, req.Body.Hash, req.Body.Chain)

	if err != nil {
		s.log.WithError(err).Debugf("Request failed")

		switch errors.Cause(err) {
		case relayer.ErrEntryNotFound, relayer.ErrChainNotFound:
			return nil, status.Errorf(codes.NotFound, err.Error())
		case relayer.ErrAlreadySubmitted:
			return nil, status.Errorf(
				codes.InvalidArgument,
				"can not resubmit state transition tx for state: %s on chain: %s", req.Body.Hash, req.Body.Chain,
			)
		default:
			s.log.WithError(err).Error("Got internal error while processing relay request")
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}

	return &types.MsgRelayResponse{Tx: tx}, nil
}

func (s *ServerImpl) StateRelays(ctx context.Context, req *types.MsgRelaysRequest) (*types.MsgRelaysResponse, error) {
	relays, err := s.relayer.StateRelays(ctx, req.Hash)
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

func (s *ServerImpl) GistRelays(ctx context.Context, req *types.MsgRelaysRequest) (*types.MsgRelaysResponse, error) {
	relays, err := s.relayer.GistRelays(ctx, req.Hash)
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
