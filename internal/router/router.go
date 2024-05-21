package router

import (
	"github.com/gorilla/mux"

	"goCalculatorYL/internal/handler"
	"goCalculatorYL/internal/service"
)

func SetupRouter(srv *service.Service) *mux.Router {
	srv.Logger.Debug("Setting up router...")

	r := mux.NewRouter()
	h := handler.NewHandler(srv)

	r.HandleFunc("/api/v1/calculate/", h.GetExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions", h.GetExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", h.GetExpressionByIDHandler).Methods("GET")
	r.HandleFunc("/internal/task", h.GetTaskHandler).Methods("GET")
	r.HandleFunc("/internal/task", h.ResultHandler).Methods("POST", "GET")

	return r
}
