package service

import (
	"fmt"
	"orchestrator/internal/models"
)

func NewExpression(expression string, id int, creator int64) *models.Expression {
	return &models.Expression{
		Id:         id,
		Expression: expression,
		Result:     0,
		Status:     "pending",
		Creator:    creator,
	}
}

func isValid(e *models.Expression) bool {
	return e.Status != "invalid"
}

func (s *MyService) invalidate(e *models.Expression) error {
	e.Status = "invalid"
	return s.repos.Expression.Update(e)
}

// GetExpressions выполняет получение списка выражений
func (s *MyService) GetExpressions(creator int64) []*models.Expression {
	s.Logger.Debugf("get all expressions (%d items)", len(s.expressions))
	var res []*models.Expression
	for _, exp := range s.expressions {
		if exp.Creator == creator {
			res = append(res, exp)
		}
	}
	return res
}

// GetExpressionById выполняет получение выражения по Id
func (s *MyService) GetExpressionById(id int, creator int64) (*models.Expression, bool) {
	exp, exists := s.expressions[id]
	if !exists {
		return nil, false
	} else if exp.Creator != creator {
		return nil, false
	}
	return exp, exists
}

// enqueueExpression добавляет новое выражение в очередь на выполнение
func (s *MyService) enqueueExpression(exp *models.Expression) error {
	_, exists := s.expressions[exp.Id]
	if exists {
		return fmt.Errorf("expression %d already exists", exp.Id)
	}
	s.expressions[exp.Id] = exp
	return nil
}

// AddExpression выполняет добавление вычисления арифметического выражения
func (s *MyService) AddExpression(req *NewExpressionRequest, creator int64) error {
	s.LastId++
	exp := NewExpression(req.Expression, s.LastId, creator)
	err := s.enqueueExpression(exp)
	if err != nil {
		return err
	}

	err = s.repos.Expression.Add(exp)
	if err != nil {
		return err
	}

	s.Logger.Infof("new expression (id: %s): %s", exp.Id, exp.Expression)
	if err := s.generateTasks(exp.Id); err != nil {
		return err
	}

	if s.tasks.taskCounter > 0 {
		exp.LastTask, _ = s.tasks.get(s.tasks.taskCounter)
	}

	return nil
}

// completeExpression выполняет всю логику при завершении вычисления выражения
func (s *MyService) completeExpression(exp *models.Expression) error {
	exp.Result = exp.LastTask.Result
	exp.Status = "done"

	if err := s.repos.Expression.Update(exp); err != nil {
		return err
	}

	s.clearTasks(exp.LastTask, true)
	exp.LastTask = nil
	s.Logger.Infof("expression (id: %d) done. Result: %f", exp.Id, exp.Result)

	return nil
}
