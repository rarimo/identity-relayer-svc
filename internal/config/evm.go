package config

import (
	"crypto/ecdsa"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	Name                string            `fig:"name"`
	BridgeAddress       common.Address    `fig:"bridge_address"`
	SubmitterPrivateKey *ecdsa.PrivateKey `fig:"submitter_private_key"`
	SubmitterAddress    common.Address    `fig:"-"`
	RPC                 *ethclient.Client `fig:"rpc"`
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
			With(figure.BaseHooks, figure.EthereumHooks, evmHooks).
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

var evmHooks = figure.Hooks{
	"*ethclient.Client": func(raw interface{}) (reflect.Value, error) {
		v, err := cast.ToStringE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "expected string")
		}

		client, err := ethclient.Dial(v)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to dial eth rpc")
		}

		return reflect.ValueOf(client), nil
	},
}
