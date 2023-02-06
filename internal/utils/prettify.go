package utils

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Prettify(data interface{}) string {
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal the data"))
	}

	return string(out)
}
