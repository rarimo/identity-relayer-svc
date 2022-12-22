package config

import (
	"reflect"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarimo/near-bridge-go/pkg/client"
	"gitlab.com/rarimo/near-bridge-go/pkg/types"
	"gitlab.com/rarimo/near-bridge-go/pkg/types/key"
)

type nearer struct {
	getter kv.Getter
	once   comfig.Once
}

type Nearer interface {
	Near() *Near
}

func NewNearer(getter kv.Getter) Nearer {
	return &nearer{
		getter: getter,
	}
}

type Near struct {
	RPC                 client.Client   `fig:"rpc"`
	SubmitterAddress    types.AccountID `fig:"submitter_address,required"`
	SubmitterPrivateKey key.KeyPair     `fig:"submitter_private_key,required"`
	BridgeAddress       types.AccountID `fig:"bridge_address,required"`
}

func (n *nearer) Near() *Near {
	return n.once.Do(func() interface{} {
		var cfg Near

		err := figure.
			Out(&cfg).
			With(nearHooks, figure.BaseHooks).
			From(kv.MustGetStringMap(n.getter, "near")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure config"))
		}

		return &cfg
	}).(*Near)
}

var nearHooks = figure.Hooks{
	"key.KeyPair": func(raw interface{}) (reflect.Value, error) {
		v, err := cast.ToStringE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected string")
		}

		keyPair, err := key.NewBase58KeyPair(v)
		if err != nil {
			panic(errors.Wrap(err, "expected a base58-encoded key pair"))
		}

		return reflect.ValueOf(keyPair), nil
	},
	"client.Client": func(raw interface{}) (reflect.Value, error) {
		v, err := cast.ToStringE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected string")
		}

		rpc, err := client.NewClient(v)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rpc client"))
		}

		return reflect.ValueOf(rpc), nil
	},
}
