package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

/**
Based on https://github.com/maticnetwork/bor/tree/ad69ccd0ba6aac4a690e6b4778987242609f4845
*/

const (
	baseFeeChangeDenominatorPreDelhi  = 8  // Bounds the amount the base fee can change between blocks before Delhi Hard Fork.
	baseFeeChangeDenominatorPostDelhi = 16 // Bounds the amount the base fee can change between blocks after Delhi Hard Fork.
)

var maticMainnetChainID = big.NewInt(137)
var mumbaiChainID = big.NewInt(80001)
var mumbaiDelhiBlock = big.NewInt(29638656)
var mainnetDelhiBlock *big.Int

func calculatePolygonBaseFee(config *params.ChainConfig, parent *types.Header) *big.Int {
	// If the current block is the first EIP-1559 block, return the InitialBaseFee.
	if !config.IsLondon(parent.Number) {
		return new(big.Int).SetUint64(params.InitialBaseFee)
	}

	var (
		parentGasTarget                = parent.GasLimit / params.ElasticityMultiplier
		parentGasTargetBig             = new(big.Int).SetUint64(parentGasTarget)
		baseFeeChangeDenominatorUint64 = baseFeeChangeDenominator(config, parent.Number)
		baseFeeChangeDenominator       = new(big.Int).SetUint64(baseFeeChangeDenominatorUint64)
	)
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.GasUsed == parentGasTarget {
		return new(big.Int).Set(parent.BaseFee)
	}
	if parent.GasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		gasUsedDelta := new(big.Int).SetUint64(parent.GasUsed - parentGasTarget)
		x := new(big.Int).Mul(parent.BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, baseFeeChangeDenominator),
			common.Big1,
		)

		return x.Add(parent.BaseFee, baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parent.GasUsed)
		x := new(big.Int).Mul(parent.BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

		return math.BigMax(
			x.Sub(parent.BaseFee, baseFeeDelta),
			common.Big0,
		)
	}
}

func baseFeeChangeDenominator(config *params.ChainConfig, number *big.Int) uint64 {
	if isDelhi(config, number) {
		return baseFeeChangeDenominatorPostDelhi
	} else {
		return baseFeeChangeDenominatorPreDelhi
	}
}

func isDelhi(config *params.ChainConfig, number *big.Int) bool {
	if config.ChainID == mumbaiChainID {
		return number.Cmp(mumbaiDelhiBlock) >= 0
	} else if config.ChainID == maticMainnetChainID {
		return number.Cmp(mainnetDelhiBlock) >= 0
	} else {
		return false
	}
}
