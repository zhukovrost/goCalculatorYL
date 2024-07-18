package service

import (
	"fmt"
)

func NewExpression(exp *NewExpressionRequest) *Expression {
	return &Expression{
		NewExpressionRequest: exp,
		Result:               0,
		Status:               "pending",
	}
}

func isValid(e *Expression) bool {
	return e.Status != "invalid"
}

// GetExpressions выполняет получение списка выражений
func (s *MyService) GetExpressions() []*Expression {
	s.Logger.Debugf("get all expressions (%d items)", len(s.expressions))
	var res []*Expression
	for _, exp := range s.expressions {
		res = append(res, exp)
	}
	return res
}

// GetExpressionById выполняет получение выражения по Id
func (s *MyService) GetExpressionById(id string) (*Expression, bool) {
	exp, exists := s.expressions[id]
	return exp, exists
}

// enqueueExpression добавляет новое выражение в очередь на выполнение
func (s *MyService) enqueueExpression(exp *Expression) error {
	_, exists := s.expressions[exp.Id]
	if exists {
		return fmt.Errorf("expression %s already exists", exp.Id)
	}
	s.expressions[exp.Id] = exp
	return nil
}

// AddExpression выполняет добавление вычисления арифметического выражения
func (s *MyService) AddExpression(req *NewExpressionRequest) error {
	exp := NewExpression(req)
	err := s.enqueueExpression(exp)
	if err != nil {
		return err
	}

	s.Logger.Infof("new expression (id: %s): %s", exp.Id, exp.Expression)
	if err := s.generateTasks(exp.Id); err != nil {
		return err
	}

	if s.tasks.taskCounter > 0 {
		exp.lastTask, _ = s.tasks.get(s.tasks.taskCounter)
	}

	return nil
}

// completeExpression выполняет всю логику при завершении вычисления выражения
func (s *MyService) completeExpression(exp *Expression) {
	exp.Result = exp.lastTask.result
	exp.Status = "done"
	s.clearTasks(exp.lastTask, true)
	exp.lastTask = nil
	s.Logger.Infof("expression (id: %s) done. result: %f", exp.Id, exp.Result)
}
