package handlers

import (
	"encoding/json"
	"net/http"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	rarimocore "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanager "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"

	"gitlab.com/rarimo/relayer-svc/internal/services"
	"gitlab.com/rarimo/relayer-svc/resources"
)

type postRelayTask struct {
	ConfirmationID string
	TransferID     string
}

func newPostRelayTaskRequest(r *http.Request) (*postRelayTask, error) {
	var request resources.RelayTaskResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	errs := ozzo.Errors{}
	res := &postRelayTask{
		ConfirmationID: request.Data.Relationships.Confirmation.Data.ID,
		TransferID:     request.Data.Relationships.Transfer.Data.ID,
	}

	errs["data.relationships.confirmation.data.id"] = ozzo.Validate(res.ConfirmationID, ozzo.Required, hexValidator)
	errs["data.relationships.transfer.data.id"] = ozzo.Validate(res.TransferID, ozzo.Required, hexValidator)

	return res, errs.Filter()
}

func PostRelayTask(w http.ResponseWriter, r *http.Request) {
	// TODO: restrict access

	request, err := newPostRelayTaskRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tm := tokenmanager.NewQueryClient(Config(r).Cosmos())
	core := rarimocore.NewQueryClient(Config(r).Cosmos())

	networks, err := tm.Params(r.Context(), new(tokenmanager.QueryParamsRequest))
	if err != nil {
		panic(errors.Wrap(err, "error getting network info"))
	}

	confirmation, err := core.Confirmation(r.Context(), &rarimocore.QueryGetConfirmationRequest{
		Root: request.ConfirmationID,
	})

	if err != nil {
		panic(errors.Wrap(err, "failed to fetch the confirmation"))
	}
	if confirmation == nil {
		ape.RenderErr(w, problems.NotFound())
		Log(r).WithField("confirmation_id", request.ConfirmationID).Error("confirmation not found")
		return
	}

	scheduler := services.NewScheduler(Config(r))
	if err := scheduler.ScheduleRelays(r.Context(), networks, confirmation.Confirmation, []string{request.TransferID}); err != nil {
		panic(errors.Wrap(err, "failed to schedule the transfers for relay"))
	}

	w.WriteHeader(http.StatusAccepted)
}
