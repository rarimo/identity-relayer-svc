package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"gitlab.com/rarimo/relayer-svc/internal/data/core"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger"
	"gitlab.com/rarimo/relayer-svc/internal/services/bridger/evm"
	"gitlab.com/rarimo/relayer-svc/internal/types"
	"gitlab.com/rarimo/relayer-svc/resources"
)

type postFeeEstimate struct {
	ConfirmationID string
	TransferID     string
}

func newPostFeeEstimate(r *http.Request) (*postFeeEstimate, error) {
	var request resources.FeeEstimateRequestResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	errs := ozzo.Errors{}
	res := &postFeeEstimate{
		ConfirmationID: request.Data.Relationships.Confirmation.Data.ID,
		TransferID:     request.Data.Relationships.Transfer.Data.ID,
	}
	errs["data.relationships.confirmation.data.id"] = ozzo.Validate(res.ConfirmationID, ozzo.Required, hexValidator)
	errs["data.relationships.transfer.data.id"] = ozzo.Validate(res.TransferID, ozzo.Required, hexValidator)

	return res, errs.Filter()
}

func PostFeeEstimate(w http.ResponseWriter, r *http.Request) {
	// TODO: add rate limits to prevent abuse
	request, err := newPostFeeEstimate(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	coreClient := core.NewCore(Config(r))
	transfer, err := coreClient.GetTransfer(r.Context(), request.ConfirmationID, request.TransferID)
	if err != nil {
		panic(errors.Wrap(err, "failed to get the transfer"))
	}
	if !types.IsEVM(transfer.Transfer.ToChain) {
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Status: "501",
			Title:  "Unsupported chain",
			Detail: fmt.Sprintf("Chain %s is not supported", transfer.Transfer.ToChain),
		})
		return
	}

	bridge := evm.NewEVMBridger(Config(r))
	feeEstimate, err := bridge.EstimateRelayFee(r.Context(), *transfer)
	if errors.Cause(err) == bridger.ErrAlreadyWithdrawn {
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Status: "409",
			Title:  "Already withdrawn",
		})
		return
	}
	if err != nil {
		panic(errors.Wrap(err, "failed to estimate gas"))
	}

	response := resources.FeeEstimateResponse{
		Data: resources.FeeEstimate{
			Key: resources.Key{
				ID:   fmt.Sprintf("%s:%d", request.TransferID, feeEstimate.CreatedAt.Unix()),
				Type: resources.FEE_ESTIMATES,
			},
			Attributes: resources.FeeEstimateAttributes{
				FeeAmount:       hexutil.EncodeBig(feeEstimate.FeeAmount),
				FeeToken:        feeEstimate.FeeToken,
				FeeTokenAddress: feeEstimate.FeeTokenAddress,
				GasEstimate:     hexutil.EncodeBig(feeEstimate.GasEstimate),
				GasToken:        feeEstimate.GasToken,
				FromChain:       feeEstimate.FromChain,
				ToChain:         feeEstimate.ToChain,
				CreatedAt:       feeEstimate.CreatedAt,
				ExpiresAt:       feeEstimate.ExpiresAt,
			},
			Relationships: resources.FeeEstimateRelationships{
				Transfer: resources.Relation{
					Data: &resources.Key{
						ID:   request.TransferID,
						Type: resources.TRANSFERS,
					},
				},
				Confirmation: resources.Relation{
					Data: &resources.Key{
						ID:   request.ConfirmationID,
						Type: resources.CONFIRMATIONS,
					},
				},
			},
		},
	}

	ape.Render(w, response)
}
