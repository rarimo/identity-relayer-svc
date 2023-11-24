package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type RelayConfig struct {
	CatchupDisabled   bool     `fig:"catchup_disabled"`
	IssuerID          []string `fig:"issuer_id"`
	DisableFiltration bool     `fig:"disable_filtration"`
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

		return cfg
	}).(RelayConfig)
}
