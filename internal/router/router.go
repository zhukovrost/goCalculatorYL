package router

import (
	"github.com/gorilla/mux"

	"goCalculatorYL/internal/handlers"
	"goCalculatorYL/internal/service"
)

func SetupRouter(srv *service.Service) *mux.Router {
	srv.Logger.Debug("Setting up router...")

	r := mux.NewRouter()
	h := handlers.NewHandler(srv)

	r.HandleFunc("/api/v1/calculate", h.AddExpressionHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", h.GetExpressionsHandler)
	r.HandleFunc("/api/v1/expressions/{id}", h.GetExpressionByIDHandler)
	r.HandleFunc("/internal/task", h.GetTaskHandler)
	r.HandleFunc("/internal/task", h.ResultHandler)

	return r
}
