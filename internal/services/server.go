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
}

func NewServer(cfg config.Config) *ServerImpl {
	return &ServerImpl{
		log:      cfg.Log(),
		listener: cfg.Listener(),
		relayer:  relayer.NewService(cfg),
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
		default:
			s.log.WithError(err).Error("got internal error while processing relay request")
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}

	return &types.MsgRelayResponse{Tx: tx}, nil
}
