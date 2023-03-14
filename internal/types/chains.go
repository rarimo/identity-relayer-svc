package types

import "golang.org/x/exp/slices"

const (
	EthereumMainnet = "Ethereum"
	Goerli          = "Goerli"
	Sepolia         = "Sepolia"

	MaticMainnet = "Matic"
	Mumbai       = "Mumbai"

	BSCMainnet = "BSC"
	Chapel     = "Chapel"

	AvalancheMainnet = "Avalanche"
	Fuji             = "Fuji"

	Solana = "Solana"
	Near   = "Near"
)

var Chains = []string{
	EthereumMainnet, Goerli, Sepolia, MaticMainnet, Mumbai, BSCMainnet, Chapel, AvalancheMainnet, Fuji, Solana, Near,
}

var evmChains = []string{
	EthereumMainnet, Goerli, Sepolia, MaticMainnet, Mumbai, BSCMainnet, Chapel, AvalancheMainnet, Fuji,
}

func IsEVM(chain string) bool {
	return slices.Contains(evmChains, chain)
}
