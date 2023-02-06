/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type FeeEstimateAttributes struct {
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	FeeAmount       string    `json:"fee_amount"`
	FeeToken        string    `json:"fee_token"`
	FeeTokenAddress string    `json:"fee_token_address"`
	FromChain       string    `json:"from_chain"`
	GasEstimate     string    `json:"gas_estimate"`
	GasToken        string    `json:"gas_token"`
	ToChain         string    `json:"to_chain"`
}
