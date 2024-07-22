package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"orchestrator/internal/service"
	"orchestrator/pkg/util"
	"strings"
)

type OrchestratorHandler struct {
	srv *service.MyService
}

func New(srv *service.MyService) *OrchestratorHandler {
	srv.Logger.Debug("Setting up orchestrator handlers...")
	return &OrchestratorHandler{
		srv: srv,
	}
}

// AddExpressionHandler выполняет добавление вычисления арифметического выражения
func (h *OrchestratorHandler) AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Debug("new POST request")

	var calculationRequest service.NewExpressionRequest
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

	calculationRequest.Id = strings.TrimSpace(calculationRequest.Id)
	if calculationRequest.Id == "" {
		calculationRequest.Id = util.GenerateId()
	}

	creator := int64(r.Context().Value("user").(float64))

	if err = h.srv.AddExpression(&calculationRequest, creator); err != nil {
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
	creator := int64(r.Context().Value("user").(float64))
	expressions := h.srv.GetExpressions(creator)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(expressions); err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}
	h.srv.Logger.Debug("successful response (200)")
}

// GetExpressionByIdHandler выполняет получение выражения по Id
func (h *OrchestratorHandler) GetExpressionByIdHandler(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Println("new GET request")
	id := chi.URLParam(r, "id")

	creator := int64(r.Context().Value("user").(float64))
	expression, exists := h.srv.GetExpressionById(id, creator)
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

// GetTaskHandler выполняет получение задачи
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

	_, err = w.Write(task)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error("failed to write response: " + err.Error())
	}
}

// SetResultHandler выполняет прием результата обработки задачи
func (h *OrchestratorHandler) SetResultHandler(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Debug("new POST request")

	var result service.CalculationResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}

	if err = h.srv.SetTaskResult(result.Id, result.Result); err != nil {
		http.Error(w, err.Error(), 404)
		h.srv.Logger.Error(err.Error())
		return
	}

	w.WriteHeader(200)
	h.srv.Logger.Debug("successful response (200)")
}

func (h *OrchestratorHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Debug("new POST request")

	var input service.UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}

	err = h.srv.Register(input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}

	w.WriteHeader(200)
	h.srv.Logger.Debug("successful response (200)")
}

func (h *OrchestratorHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.srv.Logger.Debug("new POST request")

	var input service.UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error(err.Error())
		return
	}

	token, err := h.srv.Login(input)

	if err != nil {
		switch {
		case errors.Is(err, service.InvalidCreditsError):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		h.srv.Logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(200)

	if _, err = w.Write([]byte(token + "\n")); err != nil {
		http.Error(w, err.Error(), 500)
		h.srv.Logger.Error("failed to write response: " + err.Error())
	}
}
