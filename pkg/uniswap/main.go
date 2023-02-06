package uniswap

//go:generate abigen --abi ./abi/uniswap_v2_pool.abi --pkg uniswap --type UniswapV2Pool --out ./pool_v2_generated.go
//go:generate abigen --abi ./abi/uniswap_v3_pool.abi --pkg uniswap --type UniswapV3Pool --out ./pool_v3_generated.go
//go:generate abigen --abi ./abi/erc20.abi --pkg uniswap --type ERC20 --out ./erc20_generated.go

import (
	"context"
	"math/big"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Pool interface {
	// Convert the amount of one token to another
	Convert(ctx context.Context, from, to common.Address, amount *big.Int) (*big.Int, error)
}

func getToken(rpc *ethclient.Client, address common.Address) (Token, error) {
	erc20, err := NewERC20(address, rpc)
	if err != nil {
		return Token{}, errors.Wrap(err, "failed to create ERC20")
	}

	decimals, err := erc20.Decimals(nil)
	if err != nil {
		return Token{}, errors.Wrap(err, "failed to get decimals")
	}

	return newToken(address, decimals), nil
}

func mul(x, y *big.Int) *big.Int {
	return big.NewInt(0).Mul(x, y)
}

func convert(amount, rate, denominator *big.Int) *big.Int {
	tmp := big.NewInt(0).Mul(amount, rate)

	return big.NewInt(0).Div(tmp, denominator)
}
