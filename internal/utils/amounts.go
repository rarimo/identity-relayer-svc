package utils

import (
	"errors"
	"math/big"
)

func GetAmountOrDefault(rawAmount string, defaultAmount *big.Int) (*big.Int, error) {
	if rawAmount == "" {
		return defaultAmount, nil
	}

	amount, ok := new(big.Int).SetString(rawAmount, 10)
	if !ok {
		return defaultAmount, errors.New("failed to parse amount")
	}

	return amount, nil
}
