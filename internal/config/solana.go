package config

import (
	"context"
	"reflect"

	"github.com/mr-tron/base58"
	"github.com/olegfomenko/solana-go"
	"github.com/olegfomenko/solana-go/rpc"
	"github.com/olegfomenko/solana-go/rpc/ws"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarimo/relayer-svc/internal/utils"
)

type Solaner interface {
	Solana() *Solana
}

type Solana struct {
	RPC                 *rpc.Client       `fig:"rpc,required"`
	WS                  *ws.Client        `fig:"ws,required"`
	BridgeAdmin         solana.PublicKey  `fig:"-"`
	BridgeAdminSeed     [32]byte          `fig:"bridge_admin_seed,required"`
	BridgeProgramID     solana.PublicKey  `fig:"bridge_program_id,required"`
	SubmitterPrivateKey solana.PrivateKey `fig:"submitter_private_key,required"`
}

type solaner struct {
	getter kv.Getter
	once   comfig.Once
}

func NewSolaner(getter kv.Getter) Solaner {
	return &solaner{
		getter: getter,
	}
}

func (s *solaner) Solana() *Solana {
	return s.once.Do(func() interface{} {
		config := Solana{}

		if err := figure.Out(&config).With(solanaHooks, figure.BaseHooks).From(kv.MustGetStringMap(s.getter, "solana")).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out config for solana"))
		}
		bridgeAdmin, err := solana.CreateProgramAddress([][]byte{config.BridgeAdminSeed[:]}, config.BridgeProgramID)
		if err != nil {
			panic(errors.Wrap(err, "failed to create program address"))
		}
		config.BridgeAdmin = bridgeAdmin

		return &config
	}).(*Solana)
}

var solanaHooks = figure.Hooks{
	"solana.PublicKey": func(value interface{}) (reflect.Value, error) {
		rawPubKey, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected a base58-encoded solana public key")
		}

		pubKey, err := solana.PublicKeyFromBase58(rawPubKey)
		if err != nil {
			panic(errors.Wrap(err, "valid base58-encoded solana public key expected"))
		}
		return reflect.ValueOf(pubKey), nil
	},
	"solana.PrivateKey": func(value interface{}) (reflect.Value, error) {
		rawPrivKey, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected a base58-encoded solana private key")
		}

		privKey, err := solana.PrivateKeyFromBase58(rawPrivKey)
		if err != nil {
			panic(errors.Wrap(err, "valid base58-encoded solana private key expected"))
		}
		return reflect.ValueOf(privKey), nil
	},
	"[32]uint8": func(value interface{}) (reflect.Value, error) {
		rawSeed, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected a base58-encoded solana seed")
		}

		seed, err := base58.Decode(rawSeed)
		if err != nil {
			panic(errors.Wrap(err, "valid base58-encoded solana seed expected"))
		}
		return reflect.ValueOf(utils.ToByte32(seed)), nil
	},
	"*rpc.Client": func(value interface{}) (reflect.Value, error) {
		rawRPC, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected a valid solana RPC URL")
		}

		rpcClient := rpc.New(rawRPC)
		return reflect.ValueOf(rpcClient), nil
	},
	"*ws.Client": func(value interface{}) (reflect.Value, error) {
		rawWS, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected a valid solana WS URL")
		}

		wsClient, err := ws.Connect(context.Background(), rawWS)
		if err != nil {
			panic(errors.Wrap(err, "failed to open the socket"))
		}
		return reflect.ValueOf(wsClient), nil
	},
}
