package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"goCalculatorYL/internal/config"
	"goCalculatorYL/pkg/util"
	"strconv"
)

var (
	NoTaskError = errors.New("no task")
)

type Service struct {
	Cfg         *config.Config
	Logger      *logrus.Logger
	expressions map[string]*Expression
	tasks       map[int]*Task
	taskCounter int
	lastTask    int
}

func NewService(cfg *config.Config, logger *logrus.Logger) *Service {
	return &Service{
		Cfg:         cfg,
		Logger:      logger,
		expressions: make(map[string]*Expression),
		tasks:       make(map[int]*Task),
		taskCounter: -1,
	}
}

func (s *Service) enqueueExpression(exp *Expression) error {
	_, exists := s.expressions[exp.ID]
	if exists {
		return fmt.Errorf("expression %s already exists", exp.ID)
	}
	s.expressions[exp.ID] = exp
	return nil
}

type CalculationRequest struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
}

type Expression struct {
	*CalculationRequest
	Result   float64 `json:"result"`
	Status   string  `json:"status"`
	lastTask *Task
}

func NewExpression(exp *CalculationRequest) *Expression {
	return &Expression{
		CalculationRequest: exp,
		Result:             0,
		Status:             "pending",
	}
}

func (e *Expression) IsValid() bool {
	return e.Status != "invalid"
}

// AddExpression выполняет добавление вычисления арифметического выражения
func (s *Service) AddExpression(req *CalculationRequest) error {
	exp := NewExpression(req)
	err := s.enqueueExpression(exp)
	if err != nil {
		return err
	}
	s.Logger.Infof("new expression (id: %s): %s", exp.ID, exp.Expression)
	if err := s.generateTasks(exp.ID); err != nil {
		return err
	}

	if s.taskCounter > 0 {
		exp.lastTask = s.tasks[s.taskCounter]
	}

	return nil
}

func (s *Service) completeExpression(exp *Expression) {
	exp.Result = exp.lastTask.result
	exp.Status = "done"
	s.clearTasks(exp.lastTask, true)
	exp.lastTask = nil
	s.Logger.Infof("expression (id: %s) done. result: %f", exp.ID, exp.Result)
}

func (s *Service) clearTasks(lastTask *Task, deleteCurrent bool) {
	if lastTask == nil {
		return
	}
	if isTask(lastTask.Arg1) {
		s.clearTasks(lastTask.Arg1.(*Task), true)
		lastTask.Arg1 = nil
	}
	if isTask(lastTask.Arg1) {
		s.clearTasks(lastTask.Arg1.(*Task), true)
		lastTask.Arg2 = nil
	}
	if deleteCurrent {
		delete(s.tasks, lastTask.Id)
		s.Logger.Debugf("task %d has been deleted", lastTask.Id)
	}
}

func (s *Service) generateTasks(expressionId string) error {
	s.Logger.Debug("generating tasks")

	var cnt uint
	exp := s.expressions[expressionId]
	postfix, err := util.ToPostfix(exp.Expression)

	if err != nil {
		return err
	}

	//обработка постфикса
	stack := make([]interface{}, 0) // stack to hold operands and task references
	for _, token := range postfix {
		if operand, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, operand)
		} else {
			if len(stack) < 2 {
				return fmt.Errorf("invalid postfix expression")
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			cnt++

			var task *Task
			switch token {
			case "+":
				task = s.newTask(a, b, "+", s.Cfg.Addition, expressionId)
			case "-":
				task = s.newTask(a, b, "-", s.Cfg.Subtraction, expressionId)
			case "*":
				task = s.newTask(a, b, "*", s.Cfg.Multiplication, expressionId)
			case "/":
				if !isTask(b) && b.(float64) == 0 {
					exp.Status = "invalid"
					s.Logger.Errorf("expression %v is invalid", exp.ID)
					return fmt.Errorf("division by zero")
				}
				task = s.newTask(a, b, "/", s.Cfg.Division, expressionId)
			default:
				exp.Status = "invalid"
				return fmt.Errorf("invalid operator: %s", token)
			}
			stack = append(stack, task)
		}
	}

	if len(stack) != 1 {
		return fmt.Errorf("invalid postfix expression")
	}
	s.Logger.Debugf("successfully created %d tasks", cnt)
	return nil
}

// GetExpressions выполняет получение списка выражений
func (s *Service) GetExpressions() []*Expression {
	s.Logger.Debugf("get all expressions (%d items)", len(s.expressions))
	var res []*Expression
	for _, exp := range s.expressions {
		res = append(res, exp)
	}
	return res
}

// GetExpressionByID выполняет получение списка выражений
func (s *Service) GetExpressionByID(id string) (*Expression, bool) {
	exp, exists := s.expressions[id]
	return exp, exists
}

type Task struct {
	Id            int
	Arg1          interface{}
	Arg2          interface{}
	Operation     string
	OperationTime uint
	result        float64
	expressionId  string
	isDone        bool
}

func (s *Service) newTask(arg1, arg2 interface{}, operation string, operationTime uint, expressionId string) *Task {
	s.taskCounter++
	task := &Task{
		Id:            s.taskCounter,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operation,
		OperationTime: operationTime,
		expressionId:  expressionId,
		result:        0,
		isDone:        false,
	}
	s.tasks[task.Id] = task
	return task
}

func isTask(arg interface{}) bool {
	_, ok := arg.(*Task)
	return ok
}

// GetTask выполняет получение списка выражений
func (s *Service) GetTask() (*Task, error) {
	if s.lastTask == 0 {
		task, exists := s.tasks[s.lastTask]
		if exists {
			s.lastTask++
			return task, nil
		}
		return nil, NoTaskError
	}

	task, exists := s.tasks[s.lastTask]
	if !exists {
		return nil, NoTaskError
	}

	exp := s.expressions[task.expressionId]

	if isTask(task.Arg1) && !task.Arg1.(*Task).isDone {
		return nil, NoTaskError
	}

	if isTask(task.Arg2) && !task.Arg2.(*Task).isDone {
		return nil, NoTaskError
	}

	defer func() { s.lastTask++ }()

	if !exp.IsValid() {
		return s.GetTask()
	}

	if isTask(task.Arg2) && task.Arg2.(*Task).result == 0 && task.Operation == "/" {
		exp.Status = "invalid"
		s.Logger.Errorf("expression %v is invalid", exp.ID)
		return s.GetTask()
	}

	exp.Status = "calculating"

	return task, nil
}

// SetResult выполняет прием результата обработки данных
func (s *Service) SetResult(id int, result float64) error {
	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("expression %d not found. probably, the expression is invalid", id)
	}
	task.result = result
	task.isDone = true

	exp := s.expressions[task.expressionId]
	if lastTaskId := exp.lastTask.Id; lastTaskId == task.Id {
		s.completeExpression(exp)
	}
	s.Logger.Infof("task (id: %d) done. result: %f", id, result)
	return nil
}

type Response struct {
	Task *TaskResponse `json:"task"`
}

type TaskResponse struct {
	Id            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime uint    `json:"operation_time"`
}

func (s *Service) fillResponse(task *Task) *TaskResponse {
	var arg1, arg2 float64

	if isTask(task.Arg1) {
		arg1 = task.Arg1.(*Task).result
	} else {
		arg1 = task.Arg1.(float64)
	}

	if isTask(task.Arg2) {
		arg2 = task.Arg2.(*Task).result
	} else {
		arg2 = task.Arg2.(float64)
	}

	s.clearTasks(task, false)

	return &TaskResponse{
		Id:            task.Id,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     task.Operation,
		OperationTime: task.OperationTime,
	}
}

func (s *Service) GetJSONResponse(t *Task) ([]byte, error) {
	resp := &Response{
		Task: s.fillResponse(t),
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("не получилось получить json")
	}
	return jsonData, nil
}
