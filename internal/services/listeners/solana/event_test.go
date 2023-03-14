package solana_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"gitlab.com/rarimo/relayer-svc/internal/services/listeners/solana"
)

func TestFindFeePaidEvent(t *testing.T) {
	logs := []string{
		"Program Hiso3nYYheDQFdBC2FkR3LuFiU7N1pFmvcX73mAqDLsP invoke [1]",
		"Program log: Instruction: PayFee",
		"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]",
		"Program log: Instruction: Transfer",
		"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 4645 of 191134 compute units",
		"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
		"Program data: nww01PkkGBIBAgMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMCAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAADzSq8UXLRXhVfYKdkfS9P5eEpbhRaHO8ZyEdiRWLuyNw==",
		"Program Hiso3nYYheDQFdBC2FkR3LuFiU7N1pFmvcX73mAqDLsP consumed 15274 of 200000 compute units",
		"Program Hiso3nYYheDQFdBC2FkR3LuFiU7N1pFmvcX73mAqDLsP success",
	}

	result, err := solana.FindFeePaidEvent(logs)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Amount)
	assert.Equal(t, "HNi83KdgGRQA47KWocpkZcTYGTddd2cqtNVLTgk3n8R4", result.FeeTokenMint.String())
	assert.Equal(
		t,
		"0x0102030000000000000000000000000000000000000000000000000000000000",
		hexutil.Encode(result.ConfirmationId[:]),
	)
	assert.Equal(
		t,
		"0x0302010000000000000000000000000000000000000000000000000000000000",
		hexutil.Encode(result.TransferId[:]),
	)
}
