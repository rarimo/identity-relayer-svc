package handlers

import (
	"encoding/json"
	"net/http"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

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
	errs["data.relationships.confirmation.data.id"] = ozzo.Validate(
		request.Data.Relationships.Confirmation.Data.ID,
		ozzo.Required, hexValidator,
	)
	errs["data.relationships.transfer.data.id"] = ozzo.Validate(
		request.Data.Relationships.Transfer.Data.ID,
		ozzo.Required, hexValidator,
	)
	if errs.Filter() != nil {
		return nil, errs.Filter()
	}

	res := &postRelayTask{
		ConfirmationID: request.Data.Relationships.Confirmation.Data.ID,
		TransferID:     request.Data.Relationships.Transfer.Data.ID,
	}

	return res, nil
}

func PostRelayTask(w http.ResponseWriter, r *http.Request) {
	// TODO: restrict access

	request, err := newPostRelayTaskRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	scheduler := services.NewScheduler(Config(r))
	if err := scheduler.ScheduleRelays(r.Context(), request.ConfirmationID, []string{request.TransferID}); err != nil {
		panic(errors.Wrap(err, "failed to schedule the transfers for relay"))
	}

	w.WriteHeader(http.StatusAccepted)
}
