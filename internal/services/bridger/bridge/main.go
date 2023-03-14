package bridge

import (
	"context"
	"errors"
	"math/big"
	"time"

	"gitlab.com/rarimo/relayer-svc/internal/data/core"
)

var ErrAlreadyWithdrawn = errors.New("already withdrawn")

type Bridger interface {
	// Withdraw the asset in target chain
	Withdraw(
		ctx context.Context,
		transfer core.TransferDetails,
	) error

	// EstimateRelayFee for the withdraw transaction and temporary store it
	EstimateRelayFee(
		ctx context.Context,
		transfer core.TransferDetails,
	) (FeeEstimate, error)
}

type FeeEstimate struct {
	TransferID      string    `json:"transfer_id"`
	FeeAmount       *big.Int  `json:"fee_amount"`
	FeeToken        string    `json:"fee_token"`
	FeeTokenAddress string    `json:"fee_token_address"`
	GasEstimate     *big.Int  `json:"gas_estimate"`
	GasToken        string    `json:"gas_token"`
	ToChain         string    `json:"to_chain"`
	FromChain       string    `json:"from_chain"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expired_at"`
}
