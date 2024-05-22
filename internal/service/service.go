package service

import (
	"github.com/sirupsen/logrus"
	"goCalculatorYL/internal/config"

	"net/http"
)

type Service struct {
	Cfg    *config.Config
	Logger *logrus.Logger
}

func New(cfg *config.Config, logger *logrus.Logger) *Service {
	return &Service{
		Cfg:    cfg,
		Logger: logger,
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
