package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"goCalculatorYL/internal/service"
	"goCalculatorYL/pkg/util"
	"net/http"
	"strings"
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
	h.srv.Logger.Debug("new POST request")

	var calculationRequest service.CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&calculationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.srv.Logger.Error(err.Error())
		return
	}

	if calculationRequest.Expression == "" {
		http.Error(w, "Missing required fields", http.StatusUnprocessableEntity)
		return
	}

	calculationRequest.ID = strings.TrimSpace(calculationRequest.ID)
	if calculationRequest.ID == "" {
		calculationRequest.ID = util.GenerateId()
	}

	if err = h.srv.AddExpression(&calculationRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		h.srv.Logger.Error(err.Error())
		return
	}

	// Формирование ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	h.srv.Logger.Debug("successful response (201)")
}

// GetExpressionsHandler выполняет получение списка всех выражений
func (h *OrchestratorHandler) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Debug("new GET request")
	expressions := h.srv.GetExpressions()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(expressions); err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}
	h.srv.Logger.Debug("successful response (200)")
}

// GetExpressionByIDHandler выполняет получение выражения по ID
func (h *OrchestratorHandler) GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Println("new GET request")
	vars := mux.Vars(r)
	id := vars["id"]

	expression, exists := h.srv.GetExpressionByID(id)
	if !exists {
		http.Error(w, "Expression not found", 404)
		h.srv.Logger.Errorf("Expression not found: %s", id)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(expression); err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}
	h.srv.Logger.Debug("successful response (200)")
}

// GetTaskHandler выполняет получение списка задач
func (h *OrchestratorHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := h.srv.GetTask()
	if err != nil {
		switch {
		case errors.Is(err, service.NoTaskError):
			http.Error(w, err.Error(), 404)
			return
		default:
			http.Error(w, err.Error(), 505)
			h.srv.Logger.Error(err.Error())
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	resp, err := task.GetJSONResponse()
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error("failed to write response: " + err.Error())
	}

	h.srv.Logger.Debugf("successful response: the task %d has been taken for calculation (200)", task.Id)
}

type calculationResult struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

// SetResultHandler выполняет прием результата обработки данных
func (h *OrchestratorHandler) SetResultHandler(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Debug("new POST request")

	var result calculationResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}

	if err = h.srv.SetResult(result.Id, result.Result); err != nil {
		http.Error(w, err.Error(), 404)
		h.srv.Logger.Error(err.Error())
		return
	}

	w.WriteHeader(200)
	h.srv.Logger.Debug("successful response (200)")
}
