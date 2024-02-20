package config

import (
	"crypto/tls"
	"time"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Cosmoser interface {
	Cosmos() *grpc.ClientConn
}

type cosmoser struct {
	getter kv.Getter
	once   comfig.Once
}

func NewCosmoser(getter kv.Getter) Cosmoser {
	return &cosmoser{
		getter: getter,
	}
}

func (c *cosmoser) Cosmos() *grpc.ClientConn {
	return c.once.Do(func() interface{} {
		var config struct {
			Addr string `fig:"addr"`
			TLS  bool   `fig:"enable_tls"`
		}

		if err := figure.Out(&config).From(kv.MustGetStringMap(c.getter, "cosmos")).Please(); err != nil {
			panic(err)
		}

		var client *grpc.ClientConn
		var err error

		connectSecurityOptions := grpc.WithInsecure()

		if config.TLS {
			tlsConfig := &tls.Config{
				MinVersion: tls.VersionTLS13,
			}

			connectSecurityOptions = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
		}

		client, err = grpc.Dial(config.Addr, connectSecurityOptions, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second, // wait time before ping if no activity
			Timeout: 20 * time.Second, // ping timeout
		}))
		if err != nil {
			panic(err)
		}

		return client
	}).(*grpc.ClientConn)
}
