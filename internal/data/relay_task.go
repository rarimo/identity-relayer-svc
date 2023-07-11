package data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/utils"
)

type RelayTask struct {
	OperationIndex string
	Signature      string
	Proof          string
	Origin         string
	MerklePath     []string

	RetriesLeft int
}

func NewRelayIdentityTransferTask(identityTransfer core.IdentityTransferDetails, maxRetries int) RelayTask {
	t := RelayTask{
		OperationIndex: identityTransfer.OpIndex,
		Proof:          hexutil.Encode(identityTransfer.Proof),
		Origin:         "", // TODO ??
		RetriesLeft:    maxRetries,
	}

	return t
}

func NewRelayTransferTask(transfer core.TransferDetails, maxRetries int) RelayTask {
	task := RelayTask{
		OperationIndex: transfer.Transfer.Origin,
		Signature:      transfer.Signature,
		Origin:         transfer.Origin,
		MerklePath:     make([]string, 0, len(transfer.MerklePath)),
		RetriesLeft:    maxRetries,
	}

	for _, hash := range transfer.MerklePath {
		task.MerklePath = append(task.MerklePath, hexutil.Encode(hash[:]))
	}

	return task
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
		path = append(path, utils.ToByte32(hexutil.MustDecode(hash)))
	}

	return path
}
