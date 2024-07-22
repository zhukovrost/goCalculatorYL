package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"orchestrator/internal/handlers"
	"orchestrator/internal/middleware"
)

func SetupRouter(h *handlers.OrchestratorHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/calculate", middleware.RequireAuthenticatedUser(h.AddExpressionHandler))
		r.Get("/expressions", middleware.RequireAuthenticatedUser(h.GetExpressionsHandler))
		r.Get("/expressions/{id}", middleware.RequireAuthenticatedUser(h.GetExpressionByIdHandler))

		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})

	r.Get("/internal/task", h.GetTaskHandler)
	r.Post("/internal/task", h.SetResultHandler)

	return middleware.RecoverPanic(middleware.Authenticate(r))
}
