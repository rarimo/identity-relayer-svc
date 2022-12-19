package relayer

import (
	"encoding/json"
	"math/big"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func getAmountOrDefault(rawAmount string, defaultAmount *big.Int) (*big.Int, error) {
	if rawAmount == "" {
		return defaultAmount, nil
	}

	amount, ok := new(big.Int).SetString(rawAmount, 10)
	if !ok {
		return defaultAmount, errors.New("failed to parse amount")
	}

	return amount, nil
}

func prettify(data interface{}) string {
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal the data"))
	}

	return string(out)
}
