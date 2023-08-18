package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type RelayConfig struct {
	TargetChain     string `fig:"target_chain,required"`
	CatchupDisabled bool   `fig:"catchup_disabled"`
}

func (c *config) Relay() RelayConfig {
	return c.relay.Do(func() interface{} {
		yamlName := "relay"

		var cfg RelayConfig

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, yamlName)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out "+yamlName))
		}

		if _, ok := c.EVM().GetChainByName(cfg.TargetChain); !ok {
			panic(errors.From(errors.New("target chain not found"), logan.F{
				"chain":        cfg.TargetChain,
				"config_entry": yamlName,
			}))
		}

		return cfg
	}).(RelayConfig)
}
