package uniswap

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	// Address of the token contract
	Address common.Address
	// Decimals of the token
	Decimals uint8
	// Denominator of the token decimals
	Denominator *big.Int
}

func newToken(address common.Address, decimals uint8) Token {
	return Token{
		Address:     address,
		Decimals:    decimals,
		Denominator: big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil),
	}
}
