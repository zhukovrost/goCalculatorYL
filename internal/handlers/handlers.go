package handlers

import (
	"goCalculatorYL/internal/service"

	"net/http"
)

type OrchestratorHandler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *OrchestratorHandler {
	srv.Logger.Debug("Setting up orchestrator handlers...")
	return &OrchestratorHandler{
		srv: srv,
	}
}

// AddExpressionHandler выполняет добавление вычисления арифметического выражения
func (h *OrchestratorHandler) AddExpressionHandler(w http.ResponseWriter, r *http.Request) {

}

// GetExpressionsHandler выполняет получение списка выражений
func (h *OrchestratorHandler) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {

}

// GetExpressionByIDHandler выполняет получение списка выражений
func (h *OrchestratorHandler) GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {

}

// GetTaskHandler выполняет получение списка выражений
func (h *OrchestratorHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {

}

// ResultHandler выполняет прием результата обработки данных
func (h *OrchestratorHandler) ResultHandler(w http.ResponseWriter, r *http.Request) {

}
