package config

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type EVMer interface {
	EVM() *EVM
}

type evmer struct {
	getter kv.Getter
	once   comfig.Once
}

type EVM struct {
	Chains []EVMChain `fig:"chains"`
}

type EVMChain struct {
	Name                string            `fig:"name,required"`
	BridgeAddress       common.Address    `fig:"bridge_address,required"`
	SubmitterPrivateKey *ecdsa.PrivateKey `fig:"submitter_private_key,required"`
	SubmitterAddress    common.Address    `fig:"-"`
	RPC                 *ethclient.Client `fig:"-"`
	RPCURL              string            `fig:"rpc,required"`
	ChainID             *big.Int          `fig:"-"`
}

func NewEVMer(getter kv.Getter) EVMer {
	return &evmer{
		getter: getter,
	}
}

func (e *evmer) EVM() *EVM {
	return e.once.Do(func() interface{} {
		var cfg EVM

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks, sliceHook).
			From(kv.MustGetStringMap(e.getter, "evm")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out evm config"))
		}

		return &cfg
	}).(*EVM)
}

func (e *EVMChain) TransactorOpts() *bind.TransactOpts {
	t, err := bind.NewKeyedTransactorWithChainID(e.SubmitterPrivateKey, e.ChainID)
	if err != nil {
		panic(errors.Wrap(err, "failed to create a bridge transactor"))
	}

	return t
}

func (e *EVM) GetChainByName(name string) (*EVMChain, bool) {
	for _, chain := range e.Chains {
		if chain.Name == name {
			return &chain, true
		}
	}

	return nil, false
}

var sliceHook = figure.Hooks{
	"[]config.EVMChain": func(value interface{}) (reflect.Value, error) {
		chains, err := parseEVMChain(value)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(chains), nil
	},
}

func parseEVMChain(value interface{}) ([]EVMChain, error) {
	rawSlice, err := cast.ToSliceE(value)
	if err != nil {
		return nil, errors.Wrap(err, "expected slice of EVMChain")
	}

	chains := make([]EVMChain, len(rawSlice))
	for idx, val := range rawSlice {
		raw, err := cast.ToStringMapE(val)
		if err != nil {
			return nil, errors.Wrap(err, "expected EVMChain to be map[string]interface{}")
		}

		var chain EVMChain
		if err = figure.Out(&chain).With(figure.BaseHooks, figure.EthereumHooks).From(raw).Please(); err != nil {
			return nil, errors.Wrap(err, "malformed EVMChain")
		}
		chain.RPC, err = ethclient.Dial(chain.RPCURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to dial eth rpc")
		}

		cID, err := chain.RPC.ChainID(context.TODO())
		if err != nil {
			panic(errors.Wrap(err, "failed to get chain id"))
		}
		chain.ChainID = cID
		chain.SubmitterAddress = crypto.PubkeyToAddress(chain.SubmitterPrivateKey.PublicKey)

		chains[idx] = chain
	}

	return chains, nil
}
