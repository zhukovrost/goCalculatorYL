package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"goCalculatorYL/internal/config"
	"net/http"
)

type Service struct {
	Cfg                *config.Config
	Logger             *logrus.Logger
	pendingExpressions map[string]*Expression
}

func New(cfg *config.Config, logger *logrus.Logger) *Service {
	return &Service{
		Cfg:                cfg,
		Logger:             logger,
		pendingExpressions: make(map[string]*Expression),
	}
}

func (s *Service) addExpression(exp *Expression) error {
	_, exists := s.pendingExpressions[exp.ID]
	if exists {
		return fmt.Errorf("expression %s already exists", exp.ID)
	}
	s.pendingExpressions[exp.ID] = exp
	return nil
}

type CalculationRequest struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
}

type Expression struct {
	*CalculationRequest
	Result float64 `json:"result"`
}

func NewExpression(exp *CalculationRequest) *Expression {
	return &Expression{
		CalculationRequest: exp,
		Result:             0,
	}
}

// AddExpression выполняет добавление вычисления арифметического выражения
func (s *Service) AddExpression(req *CalculationRequest) error {
	exp := NewExpression(req)
	err := s.addExpression(exp)
	if err != nil {
		return err
	}
	s.Logger.Infof("new expression (id: %s): %s", exp.ID, exp.Expression)
	return nil
}

// GetExpressions выполняет получение списка выражений
func (s *Service) GetExpressions() []*Expression {
	s.Logger.Debugf("get all expressions (%d items)", len(s.pendingExpressions))
	var res []*Expression
	for _, exp := range s.pendingExpressions {
		res = append(res, exp)
	}
	return res
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
