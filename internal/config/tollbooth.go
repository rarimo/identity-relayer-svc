package config

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/olegfomenko/solana-go"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Tollboother interface {
	Tollbooth() *Tollbooth
}

type tollboother struct {
	getter kv.Getter
	once   comfig.Once
}

type Tollbooth struct {
	FeeTokenTicker string `fig:"fee_token_ticker,required"`

	// Chain specific configs
	Solana           SolanaTollboothConfig `fig:"solana,required"`
	EVM              []EVMTollboothConfig  `fig:"evm,required"`
	evmConfigsByName map[string]EVMTollboothConfig
}

type SolanaTollboothConfig struct {
	GasTokenTicker   string           `fig:"-"`
	FeeTokenMint     solana.PublicKey `fig:"fee_token_mint,required"`
	TollboothAccount solana.PublicKey `fig:"tollbooth_account,required"`
	RelayFee         *big.Int         `fig:"relay_fee,required"`
}

type EVMTollboothConfig struct {
	Chain                  string         `fig:"chain,required"`
	GasTokenTicker         string         `fig:"gas_token_ticker,required"`
	TollboothAddress       common.Address `fig:"tollbooth_address,required"`
	WrappedGasTokenAddress common.Address `fig:"gas_token_address,required"`
	FeeTokenAddress        common.Address `fig:"fee_token_address,required"`
	UniswapPoolAddress     common.Address `fig:"uniswap_pool_address,required"`
	UseUniswapV2           bool           `fig:"use_uniswap_v2"`
}

func NewTollboother(getter kv.Getter) Tollboother {
	return &tollboother{
		getter: getter,
	}
}

func (t *tollboother) Tollbooth() *Tollbooth {
	return t.once.Do(func() interface{} {
		var tollbooth Tollbooth
		err := figure.
			Out(&tollbooth).
			With(figure.BaseHooks, tollboothSliceHook, solanaHooks).
			From(kv.MustGetStringMap(t.getter, "tollbooth")).
			Please()
		if err != nil {
			panic(err)
		}

		tollbooth.evmConfigsByName = make(map[string]EVMTollboothConfig)
		for _, chain := range tollbooth.EVM {
			tollbooth.evmConfigsByName[chain.Chain] = chain
		}

		tollbooth.Solana.GasTokenTicker = "SOL"

		return &tollbooth
	}).(*Tollbooth)
}

func (t *Tollbooth) GetEVMConfig(chain string) EVMTollboothConfig {
	return t.evmConfigsByName[chain]
}

var tollboothSliceHook = figure.Hooks{
	"[]config.EVMTollboothConfig": func(value interface{}) (reflect.Value, error) {
		chains, err := parseEVMTollbooth(value)
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(chains), nil
	},
}

func parseEVMTollbooth(value interface{}) ([]EVMTollboothConfig, error) {
	rawSlice, err := cast.ToSliceE(value)
	if err != nil {
		return nil, errors.Wrap(err, "expected slice of EVMTollboothConfig")
	}

	chains := make([]EVMTollboothConfig, len(rawSlice))
	for idx, val := range rawSlice {
		raw, err := cast.ToStringMapE(val)
		if err != nil {
			return nil, errors.Wrap(err, "expected EVMTollboothConfig to be map[string]interface{}")
		}

		var chain EVMTollboothConfig
		if err = figure.Out(&chain).With(figure.BaseHooks, figure.EthereumHooks).From(raw).Please(); err != nil {
			return nil, errors.Wrap(err, "malformed EVMTollboothConfig")
		}
		chains[idx] = chain
	}

	return chains, nil
}
