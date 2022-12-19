package api

import (
	"context"
	"net"
	"time"

	"gitlab.com/rarify-protocol/relayer-svc/internal/config"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type api struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
	cfg      config.Config
}

func Run(cfg config.Config, ctx context.Context) {
	s := &api{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
		cfg:      cfg,
	}
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		panic(errors.Wrap(err, "cop failed"))
	}

	s.log.WithFields(logan.F{
		"service": "api",
		"addr":    cfg.Listener().Addr(),
	}).Info("listening for http requests")
	ape.Serve(ctx, r, cfg, ape.ServeOpts{
		WriteTimeout: time.Second * 30,
	})
}
