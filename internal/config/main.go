package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/rarimo/relayer-svc/internal/data/redis"
)

type Config interface {
	comfig.Logger
	redis.Rediserer
	Tenderminter
	Cosmoser
	EVMer

	Relay() RelayConfig
}

type config struct {
	comfig.Logger
	redis.Rediserer
	getter kv.Getter
	Tenderminter
	Cosmoser
	EVMer

	relay comfig.Once
}

func New(getter kv.Getter) Config {
	logger := comfig.NewLogger(getter, comfig.LoggerOpts{})
	return &config{
		getter:       getter,
		Logger:       logger,
		Rediserer:    redis.NewRediserer(getter, logger.Log()),
		Tenderminter: NewTenderminter(getter),
		Cosmoser:     NewCosmoser(getter),
		EVMer:        NewEVMer(getter),
	}
}
