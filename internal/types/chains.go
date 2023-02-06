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

	SolanaMainnet = "Solana"
	NearMainnet   = "Near"
)

var Chains = []string{
	EthereumMainnet, Goerli, Sepolia, MaticMainnet, Mumbai, BSCMainnet, Chapel, AvalancheMainnet, Fuji, SolanaMainnet, NearMainnet,
}

func IsEVM(chain string) bool {
	return slices.Contains(Chains, chain)
}
