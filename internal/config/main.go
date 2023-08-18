package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	Tenderminter
	Cosmoser
	EVMer

	Relay() RelayConfig
}

type config struct {
	comfig.Logger
	pgdb.Databaser

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
		Databaser:    pgdb.NewDatabaser(getter),
		Tenderminter: NewTenderminter(getter),
		Cosmoser:     NewCosmoser(getter),
		EVMer:        NewEVMer(getter),
	}
}
