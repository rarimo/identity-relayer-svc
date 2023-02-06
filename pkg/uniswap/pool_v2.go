package uniswap

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var q112 = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(112), nil)
var v2TWAPIntervalBlocks = uint64(100)

type poolV2 struct {
	logger *logan.Entry
	rpc    *ethclient.Client
	pool   *UniswapV2Pool
	token0 Token
	token1 Token
}

func NewPoolV2(logger *logan.Entry, rpc *ethclient.Client, poolAddress common.Address) (Pool, error) {
	poolContract, err := NewUniswapV2Pool(poolAddress, rpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create UniswapV3Pool")
	}

	p := poolV2{
		logger: logger,
		rpc:    rpc,
		pool:   poolContract,
	}

	token0, err := poolContract.UniswapV2PoolCaller.Token0(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token0 address")
	}
	p.token0, err = getToken(p.rpc, token0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token0")
	}

	token1, err := poolContract.UniswapV2PoolCaller.Token1(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token1 address")
	}
	p.token1, err = getToken(p.rpc, token1)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token1")
	}

	return &p, nil
}

func (p *poolV2) Convert(ctx context.Context, from, to common.Address, amount *big.Int) (*big.Int, error) {
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
		return convert(amount, token1, q112), nil
	}

	return convert(amount, token0, q112), nil
}

func (p *poolV2) getTWAP(ctx context.Context) (token0 *big.Int, token1 *big.Int, err error) {
	blockNumber, err := p.rpc.BlockNumber(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get the block number")
	}

	token0Prev, token1Prev, tPrev, err := p.getAccumulatedPrice(ctx, blockNumber-v2TWAPIntervalBlocks)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get the accumulated price")
	}

	token0Curr, token1Curr, tCurr, err := p.getAccumulatedPrice(ctx, blockNumber)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get the accumulated price")
	}

	timePassed := big.NewInt(0).SetUint64(tCurr - tPrev)
	token0Diff := big.NewInt(0).Sub(token0Curr, token0Prev)
	token1Diff := big.NewInt(0).Sub(token1Curr, token1Prev)
	token0 = big.NewInt(0).Div(token0Diff, timePassed)
	token1 = big.NewInt(0).Div(token1Diff, timePassed)

	return token0, token1, nil
}

func (p *poolV2) getAccumulatedPrice(
	ctx context.Context,
	blockNumber uint64,
) (token0, token1 *big.Int, blockTimestamp uint64, err error) {
	bn := big.NewInt(0).SetUint64(blockNumber)
	opts := bind.CallOpts{Context: ctx, BlockNumber: bn}

	if reserves, err := p.pool.UniswapV2PoolCaller.GetReserves(&opts); err != nil {
		return nil, nil, 0, errors.Wrap(err, "failed to get the reserves")
	} else if reserves.Reserve0.Cmp(big.NewInt(0)) == 0 || reserves.Reserve1.Cmp(big.NewInt(0)) == 0 {
		return nil, nil, 0, errors.New("pool is empty")
	}

	if token0, err = p.pool.UniswapV2PoolCaller.Price0CumulativeLast(&opts); err != nil {
		return nil, nil, 0, errors.Wrap(err, "failed to get the token0 price")
	}
	if token1, err = p.pool.UniswapV2PoolCaller.Price1CumulativeLast(&opts); err != nil {
		return nil, nil, 0, errors.Wrap(err, "failed to get the token1 price")
	}

	block, err := p.rpc.HeaderByNumber(ctx, bn)
	if err != nil {
		return nil, nil, 0, errors.Wrap(err, "failed to get the block timestamp")
	}
	t := block.Time % (1 << 32)

	return token0, token1, t, nil
}
