package data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
)

type RelayTask struct {
	OperationIndex string
	Proof          string
	RetriesLeft    int
}

func NewRelayIdentityTransferTask(identityTransfer core.IdentityTransferDetails, maxRetries int) RelayTask {
	return RelayTask{
		OperationIndex: identityTransfer.OpIndex,
		Proof:          hexutil.Encode(identityTransfer.Proof),
		RetriesLeft:    maxRetries,
	}
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
