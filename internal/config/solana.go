package config

import (
	"context"

	"github.com/mr-tron/base58"
	"github.com/olegfomenko/solana-go"
	"github.com/olegfomenko/solana-go/rpc"
	"github.com/olegfomenko/solana-go/rpc/ws"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarify-protocol/relayer-svc/internal/helpers"
)

type Solaner interface {
	Solana() *Solana
}

type Solana struct {
	RPC                 *rpc.Client
	WS                  *ws.Client
	BridgeAdmin         solana.PublicKey
	BridgeAdminSeed     [32]byte
	BridgeProgramID     solana.PublicKey
	SubmitterPrivateKey solana.PrivateKey
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
		var config struct {
			RPC                 string `fig:"rpc"`
			WS                  string `fig:"ws"`
			SubmitterPrivateKey string `fig:"submitter_private_key"`
			BridgeAdminSeed     string `fig:"bridge_admin_seed"`
			BridgeProgramID     string `fig:"bridge_program_id"`
		}

		if err := figure.Out(&config).From(kv.MustGetStringMap(s.getter, "solana")).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out config for solana"))
		}

		result := Solana{}
		result.RPC = rpc.New(config.RPC)

		// Create a new WS client (used for confirming transactions)
		wsClient, err := ws.Connect(context.Background(), config.WS)
		if err != nil {
			panic(errors.Wrap(err, "failed to open the socket"))
		}
		result.WS = wsClient

		key, err := solana.PrivateKeyFromBase58(config.SubmitterPrivateKey)
		if err != nil {
			panic(errors.Wrap(err, "valid base58-encoded solana private key expected"))
		}
		result.SubmitterPrivateKey = key

		programID, err := solana.PublicKeyFromBase58(config.BridgeProgramID)
		if err != nil {
			panic(errors.Wrap(err, "valid base58-encoded solana program ID expected"))
		}
		result.BridgeProgramID = programID

		seed, err := base58.Decode(config.BridgeAdminSeed)
		if err != nil {
			panic(errors.Wrap(err, "valid base58-encoded solana seed expected"))
		}
		result.BridgeAdminSeed = helpers.ToByte32(seed)

		bridgeAdmin, err := solana.CreateProgramAddress([][]byte{seed[:]}, programID)
		if err != nil {
			panic(errors.Wrap(err, "failed to create the bridge admin address"))
		}
		result.BridgeAdmin = bridgeAdmin

		return &result
	}).(*Solana)
}
