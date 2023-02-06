package uniswap

import (
	"context"
	"math/big"

	"github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const v3TWAPIntervalSeconds = 30 * 60 // 30 minutes
var q192 = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(192), nil)

type poolV3 struct {
	logger *logan.Entry
	rpc    *ethclient.Client
	pool   *UniswapV3Pool
	token0 Token
	token1 Token
}

func (p *poolV3) Token0() Token {
	return p.token0
}

func (p *poolV3) Token1() Token {
	return p.token1
}

func NewPoolV3(logger *logan.Entry, rpc *ethclient.Client, poolAddress common.Address) (Pool, error) {
	poolContract, err := NewUniswapV3Pool(poolAddress, rpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create UniswapV3Pool")
	}

	p := poolV3{
		logger: logger,
		rpc:    rpc,
		pool:   poolContract,
	}

	token0, err := poolContract.UniswapV3PoolCaller.Token0(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token0 address")
	}
	p.token0, err = getToken(p.rpc, token0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token0")
	}

	token1, err := poolContract.UniswapV3PoolCaller.Token1(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token1 address")
	}
	p.token1, err = getToken(p.rpc, token1)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token1")
	}

	return &p, nil
}

func (p *poolV3) Convert(ctx context.Context, from, to common.Address, amount *big.Int) (*big.Int, error) {
	if from == to {
		return amount, nil
	}

	if from != p.token0.Address && from != p.token1.Address {
		return nil, errors.From(errors.New("from token is not in the pool"), logan.F{
			"from": from.Hex(),
		})
	}

	if to != p.token0.Address && to != p.token1.Address {
		return nil, errors.From(errors.New("to token is not in the pool"), logan.F{
			"to": to.Hex(),
		})
	}

	token0, token1, err := p.getTWAP(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the TWAP")
	}

	if to == p.token0.Address {
		return convert(amount, token1, p.token1.Denominator), nil
	}

	return convert(amount, token0, p.token0.Denominator), nil
}

func (p *poolV3) getTWAP(ctx context.Context) (token0 *big.Int, token1 *big.Int, err error) {
	opts := bind.CallOpts{Context: ctx}

	secondsAgo := []uint32{uint32(v3TWAPIntervalSeconds), 0}
	ticks, err := p.pool.Observe(&opts, secondsAgo)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get ticks")
	}

	tick := int((ticks.TickCumulatives[1].Int64() - ticks.TickCumulatives[0].Int64())) / v3TWAPIntervalSeconds

	sqrtPriceX96, err := utils.GetSqrtRatioAtTick(tick)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get sqrtPriceX96")
	}

	priceX192 := big.NewInt(0).Mul(sqrtPriceX96, sqrtPriceX96)
	token0 = big.NewInt(0).Div(mul(priceX192, p.token0.Denominator), q192)
	token1 = big.NewInt(0).Div(q192, mul(priceX192, p.token1.Denominator))

	return token0, token1, nil
}
