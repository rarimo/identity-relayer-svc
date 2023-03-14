package solana

import (
	"bytes"
	"encoding/base64"
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/near/borsh-go"
	"github.com/olegfomenko/solana-go"
)

// anchor ID of the FeePaid event: sha256(`event:${name}`))[0:8]
var eventDiscriminator = hexutil.MustDecode("0x9f0c34d4f9241812")

// NOTE: order of fields is important, it must match the order of fields in the program
type FeePaidEvent struct {
	ConfirmationId [32]byte
	TransferId     [32]byte
	Amount         uint64
	FeeTokenMint   solana.PublicKey
}

func FindFeePaidEvent(logs []string) (FeePaidEvent, error) {
	var rawEvent string
	for _, log := range logs {
		if strings.Contains(log, "Program data: ") {
			rawEvent = log
			break
		}
	}

	rawData, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(rawEvent, "Program data: "))
	if err != nil {
		return FeePaidEvent{}, errors.Wrap(err, "failed to decode the event data")
	}
	if len(rawData) < 8 {
		return FeePaidEvent{}, errors.New("invalid event data")
	}
	if !bytes.Equal(rawData[:8], eventDiscriminator) {
		return FeePaidEvent{}, errors.New("invalid event discriminator")
	}

	var event FeePaidEvent
	if err := borsh.Deserialize(&event, rawData[8:]); err != nil {
		return FeePaidEvent{}, errors.Wrap(err, "failed to deserialize the event data")
	}

	return event, nil
}
