package services

import (
	"context"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	client "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/operation"
	"golang.org/x/exp/slices"

	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"

	merkle "gitlab.com/rarimo/go-merkle"

	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"gitlab.com/rarimo/relayer-svc/internal/config"
	"gitlab.com/rarimo/relayer-svc/internal/data"
	"gitlab.com/rarimo/relayer-svc/internal/services/relayer"
)

const (
	DepositChanSize = 100
)

type Scheduler interface {
	ScheduleRelays(
		ctx context.Context,
		networks *tokenmanager.QueryParamsResponse,
		confirmation rarimocore.Confirmation,
		transferIndexes []string,
	) error
}

type scheduler struct {
	client       *http.HTTP
	log          *logan.Entry
	cosmos       client.ServiceClient
	tokenmanager tokenmanager.QueryClient
	rarimocore   rarimocore.QueryClient
	relayQueue   rmq.Queue
}

func NewScheduler(cfg config.Config) Scheduler {
	return newScheduler(cfg)
}

func newScheduler(cfg config.Config) *scheduler {
	return &scheduler{
		client:       cfg.Tendermint(),
		log:          cfg.Log().WithField("service", "scheduler"),
		cosmos:       client.NewServiceClient(cfg.Cosmos()),
		tokenmanager: tokenmanager.NewQueryClient(cfg.Cosmos()),
		rarimocore:   rarimocore.NewQueryClient(cfg.Cosmos()),
		relayQueue:   cfg.Redis().OpenRelayQueue(),
	}
}

func RunScheduler(cfg config.Config, ctx context.Context) {
	s := newScheduler(cfg)
	s.run(ctx)
}

func (s *scheduler) run(ctx context.Context) error {
	out, err := s.client.Subscribe(
		ctx,
		"scheduler",
		"tm.event='Tx' AND message.action='create_confirmation'",
		DepositChanSize,
	)
	s.log.Info("listening for confirmations")
	if err != nil {
		return errors.Wrap(err, "can not subscribe")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c := <-out:
			// Delay for indexing tx in core databases
			time.Sleep(time.Second * 5)
			tx, err := s.cosmos.GetTx(ctx, &client.GetTxRequest{Hash: c.Events["tx.hash"][0]})
			if err != nil {
				s.log.WithError(err).Error("error fetching tx by hash")
				continue
			}
			networks, err := s.tokenmanager.Params(ctx, new(tokenmanager.QueryParamsRequest))
			if err != nil {
				s.log.WithError(err).Error("error getting network info")
				continue
			}

			for _, message := range tx.Tx.Body.Messages {
				if message.TypeUrl != "/rarifyprotocol.rarimocore.rarimocore.MsgCreateConfirmation" {
					continue
				}

				msg := rarimocore.MsgCreateConfirmation{}
				if err = msg.Unmarshal(message.Value); err != nil {
					s.log.WithError(err).Error("failed to unmarshal message")
					continue
				}
				confirmation := rarimocore.Confirmation{
					Root:           msg.Root,
					SignatureECDSA: msg.SignatureECDSA,
					Indexes:        msg.Indexes,
					Creator:        msg.Creator,
				}

				if err := s.ScheduleRelays(ctx, networks, confirmation, msg.Indexes); err != nil {
					s.log.WithError(err).Error("failed to schedule")
				}
			}
		}
	}
}

func (s *scheduler) ScheduleRelays(
	ctx context.Context,
	networks *tokenmanager.QueryParamsResponse,
	confirmation rarimocore.Confirmation,
	transferIndexes []string,
) error {
	log := s.log.WithField("merkle_root", confirmation.Root)
	log.Info("processing a confirmation")

	tasks := []data.RelayTask{}
	operations := []*operation.TransferContent{}
	contents := []merkle.Content{}

	for _, id := range confirmation.Indexes {
		operation, err := s.rarimocore.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: id})
		if err != nil {
			log.WithError(err).Error("error getting operation entry")
			continue
		}
		transfer := rarimocore.Transfer{}
		if err := transfer.Unmarshal(operation.Operation.Details.Value); err != nil {
			log.WithError(err).Error("failed to unmarshal  transfer")
			continue
		}

		content, err := s.hashTransfer(ctx, transfer, networks.Params.Networks)
		if err != nil {
			log.WithError(err).Error("failed to hash content of transfer")
			continue
		}
		contents = append(contents, content)
		operations = append(operations, content)

		tasks = append(tasks, data.RelayTask{
			OperationIndex: id,
			Signature:      confirmation.SignatureECDSA,
			Origin:         hexutil.Encode(content.Origin[:]),
			RetriesLeft:    relayer.MaxRetries,
		})
	}

	tree := merkle.NewTree(crypto.Keccak256, contents...)
	for i, operation := range operations {
		rawPath, ok := tree.Path(operation)
		if !ok {
			panic(fmt.Errorf("failed to build Merkle tree"))
		}
		tasks[i].MerklePath = make([]string, 0, len(rawPath))
		for _, hash := range rawPath {
			tasks[i].MerklePath = append(tasks[i].MerklePath, hexutil.Encode(hash))
		}

		log.
			WithFields(logan.F{
				"op_index": tasks[i].OperationIndex,
				"op_hash":  operation.CalculateHash(),
			}).
			Infof("scheduling relay to %s", operation.TargetNetwork)
	}

	rawTasks := [][]byte{}
	for _, task := range tasks {
		if slices.Contains(transferIndexes, task.OperationIndex) {
			rawTasks = append(rawTasks, task.Marshal())
		}
	}

	if err := s.relayQueue.PublishBytes(rawTasks...); err != nil {
		return errors.Wrap(err, "failed to publish tasks")
	}

	log.Infof("scheduled %d transfers for relay", len(rawTasks))

	return nil

}

func (s *scheduler) hashTransfer(
	ctx context.Context,
	transfer rarimocore.Transfer,
	networks map[string]*tokenmanager.ChainParams,
) (*operation.TransferContent, error) {
	infoResp, err := s.tokenmanager.Info(ctx, &tokenmanager.QueryGetInfoRequest{Index: transfer.TokenIndex})
	if err != nil {
		return nil, errors.Wrap(err, "error getting token info entry")
	}

	itemResp, err := s.tokenmanager.Item(ctx, &tokenmanager.QueryGetItemRequest{
		TokenAddress: infoResp.Info.Chains[transfer.ToChain].TokenAddress,
		TokenId:      infoResp.Info.Chains[transfer.ToChain].TokenId,
		Chain:        transfer.ToChain,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting token item entry")
	}

	return pkg.GetTransferContent(&itemResp.Item, networks[transfer.ToChain], &transfer)
}
