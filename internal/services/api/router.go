package api

import (
	"gitlab.com/rarimo/relayer-svc/internal/services/api/handlers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *api) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxConfig(s.cfg),
		),
	)

	r.Route("/relayer", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/fee_estimates", handlers.PostFeeEstimate)
			r.Post("/relay_tasks", handlers.PostRelayTask)
		})
	})

	return r
}
