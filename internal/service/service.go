package service

import (
	"github.com/sirupsen/logrus"
	"goCalculatorYL/internal/config"
	"net/http"
)

type Service struct {
	Cfg                *config.Config
	Logger             *logrus.Logger
	pendingExpressions []*Expression
}

func New(cfg *config.Config, logger *logrus.Logger) *Service {
	return &Service{
		Cfg:    cfg,
		Logger: logger,
	}
}

func (s *Service) addExpression(exp *Expression) {
	s.pendingExpressions = append(s.pendingExpressions, exp)
}

type CalculationRequest struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
}

type Expression struct {
	*CalculationRequest
	result float64 `json:"result"`
}

func NewExpression(exp *CalculationRequest) *Expression {
	return &Expression{
		CalculationRequest: exp,
		result:             0,
	}
}

// AddExpression выполняет добавление вычисления арифметического выражения
func (s *Service) AddExpression(req *CalculationRequest) {
	exp := NewExpression(req)
	s.addExpression(exp)
	s.Logger.Infof("new expression (id: %s): %s", exp.ID, exp.Expression)
}

// GetExpressions выполняет получение списка выражений
func (s *Service) GetExpressions() {

}

// GetExpressionByID выполняет получение списка выражений
func (s *Service) GetExpressionByID(w http.ResponseWriter, r *http.Request) {

}

// GetTask выполняет получение списка выражений
func (s *Service) GetTask() {

}

// Result выполняет прием результата обработки данных
func (s *Service) Result() {

}
