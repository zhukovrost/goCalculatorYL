package service

import (
	"goCalculatorYL/internal/config"
	"net/http"
)

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

// AddExpression выполняет добавление вычисления арифметического выражения
func (s *Service) AddExpression() {

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
