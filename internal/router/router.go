package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"orchestrator/internal/handlers"
)

func SetupRouter(h *handlers.OrchestratorHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/calculate", h.AddExpressionHandler)
		r.Get("/expressions", h.GetExpressionsHandler)
		r.Get("/expressions/{id}", h.GetExpressionByIdHandler)
	})

	r.Get("/internal/task", h.GetTaskHandler)
	r.Post("/internal/task", h.SetResultHandler)

	return r
}
