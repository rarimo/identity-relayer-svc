package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	rediser "gitlab.com/rarimo/relayer-svc/internal/data"
)

type Config interface {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	rediser.Rediserer
	Tenderminter
	Cosmoser
	EVMer
	Solaner
	Nearer
}

type config struct {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	rediser.Rediserer
	getter kv.Getter
	Tenderminter
	Cosmoser
	EVMer
	Solaner
	Nearer
}

func New(getter kv.Getter) Config {
	logger := comfig.NewLogger(getter, comfig.LoggerOpts{})
	return &config{
		getter:       getter,
		Copuser:      copus.NewCopuser(getter),
		Listenerer:   comfig.NewListenerer(getter),
		Logger:       logger,
		Rediserer:    rediser.NewRediserer(getter, logger.Log()),
		Tenderminter: NewTenderminter(getter),
		Cosmoser:     NewCosmoser(getter),
		EVMer:        NewEVMer(getter),
		Solaner:      NewSolaner(getter),
		Nearer:       NewNearer(getter),
	}
}
