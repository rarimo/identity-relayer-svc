package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/rarimo/relayer-svc/internal/data/redis"
)

type Config interface {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	redis.Rediserer
	Tenderminter
	Cosmoser
	EVMer
}

type config struct {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	redis.Rediserer
	getter kv.Getter
	Tenderminter
	Cosmoser
	EVMer
}

func New(getter kv.Getter) Config {
	logger := comfig.NewLogger(getter, comfig.LoggerOpts{})
	return &config{
		getter:       getter,
		Copuser:      copus.NewCopuser(getter),
		Listenerer:   comfig.NewListenerer(getter),
		Logger:       logger,
		Rediserer:    redis.NewRediserer(getter, logger.Log()),
		Tenderminter: NewTenderminter(getter),
		Cosmoser:     NewCosmoser(getter),
		EVMer:        NewEVMer(getter),
	}
}
