package redis

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Rediserer interface {
	Redis() Rediser
}

type config struct {
	Addr string `fig:"addr,required"`
}

type rediserer struct {
	redisOnce comfig.Once
	getter    kv.Getter
	logger    *logan.Entry
}

func NewRediserer(getter kv.Getter, logger *logan.Entry) Rediserer {
	return &rediserer{
		getter: getter,
		logger: logger,
	}
}

func (s *rediserer) Redis() Rediser {
	return s.redisOnce.Do(func() interface{} {
		var cfg config

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(s.getter, "redis")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out redis config"))
		}

		return NewRediser(cfg, s.logger)
	}).(Rediser)
}
