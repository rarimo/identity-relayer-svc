package uniswap

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

// ETH/USDT pool
var poolV3Address = common.HexToAddress("0x4e68ccd3e89f51c3074ca5072bbac773960dfa36")
var ethereumWETHAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
var ethUSDTAddress = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
var rpcEndpoint = "https://ethereum-mainnet-rpc.allthatnode.com"

// // BNB/USDC pool
var bscPoolV2Address = common.HexToAddress("0x16b9a82891338f9bA80E2D6970FddA79D1eb0daE")
var bscWBNBAddress = common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
var bscUSDCAddress = common.HexToAddress("0x55d398326f99059fF775485246999027B3197955")
var rpcBSCEndpoint = "https://bsc-mainnet-rpc.allthatnode.com"

// AVAX/USDT pool
var avaxPoolV2Address = common.HexToAddress("0xe3bA3d5e3F98eefF5e9EDdD5Bd20E476202770da")
var avaxWAVAXAddress = common.HexToAddress("0xb31f66aa3c1e785363f0875a1b74e27b85fd66c7")
var avaxUSDTAddress = common.HexToAddress("0x9702230a8ea53601f5cd2dc00fdbc13d4df4a8c7")
var rpcAvaxEndpoint = "https://api.avax.network/ext/bc/C/rpc"

func TestV3Convert(t *testing.T) {
	client, err := ethclient.Dial(rpcEndpoint)
	assert.NoError(t, err)

	uniswap, err := NewPoolV3(nil, client, poolV3Address)
	assert.NoError(t, err)

	amount := big.NewInt(10000000000000000) // 0.01 ETH
	converted0to1, err := uniswap.Convert(context.TODO(), ethereumWETHAddress, ethUSDTAddress, amount)
	assert.NoError(t, err)

	fmt.Println("Converted:", converted0to1.Text(10))
}

func TestV2ConvertBSC(t *testing.T) {
	client, err := ethclient.Dial(rpcBSCEndpoint)
	assert.NoError(t, err)

	uniswap, err := NewPoolV2(nil, client, bscPoolV2Address)
	assert.NoError(t, err)
	amount := big.NewInt(10000000000000000) // 0.01 BNB
	converted0to1, err := uniswap.Convert(context.TODO(), bscWBNBAddress, bscUSDCAddress, amount)
	assert.NoError(t, err)

	fmt.Println("Converted:", converted0to1.Text(10))
}

func TestV2ConvertAvax(t *testing.T) {
	client, err := ethclient.Dial(rpcAvaxEndpoint)
	assert.NoError(t, err)

	uniswap, err := NewPoolV2(nil, client, avaxPoolV2Address)
	assert.NoError(t, err)
	amount := big.NewInt(10000000000000000) // 0.01 AVAX
	converted0to1, err := uniswap.Convert(context.TODO(), avaxWAVAXAddress, avaxUSDTAddress, amount)
	assert.NoError(t, err)

	fmt.Println("Converted:", converted0to1.Text(10))
}
