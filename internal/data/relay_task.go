package data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarify-protocol/relayer-svc/internal/helpers"
)

type RelayTask struct {
	OperationIndex string
	Signature      string
	Origin         string
	MerklePath     []string

	RetriesLeft int
}

func (r RelayTask) Marshal() []byte {
	marshaled, err := json.Marshal(r)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal the relay task"))
	}

	return marshaled
}

func (r *RelayTask) Unmarshal(data string) {
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal the relay task"))
	}
}

func (r RelayTask) MustParseMerklePath() [][32]byte {
	path := [][32]byte{}
	for _, hash := range r.MerklePath {
		path = append(path, helpers.ToByte32(hexutil.MustDecode(hash)))
	}

	return path
}
