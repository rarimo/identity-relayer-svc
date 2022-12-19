package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/logan/v3/errors"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

var hexValidator ozzo.Rule = ozzo.By(validateHex)

func validateHex(value interface{}) error {
	s, _ := value.(string)
	bytes, err := hexutil.Decode(s)
	if err != nil {
		return errors.Wrap(err, "invalid hex string")
	}
	if len(bytes) == 0 {
		return errors.New("empty hex string")
	}

	return nil
}
