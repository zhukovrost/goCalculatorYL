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
	taskCounter int // Переменная, для хранения id каждого новой задачи
	lastTask    int // Переменная, для хранения id последней выполненной задачи
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

// enqueueExpression добавляет новое выражение в очередь на выполнение
func (s *Service) enqueueExpression(exp *Expression) error {
	_, exists := s.expressions[exp.Id]
	if exists {
		return fmt.Errorf("expression %s already exists", exp.Id)
	}
	s.expressions[exp.Id] = exp
	return nil
}

// NewExpressionRequest является входными данными при приёме нового выражения
type NewExpressionRequest struct {
	Id         string `json:"id"`
	Expression string `json:"expression"`
}

// Expression является выражением, которое нужно вычислить
type Expression struct {
	*NewExpressionRequest
	Result   float64 `json:"result"`
	Status   string  `json:"status"`
	lastTask *Task
}

func NewExpression(exp *NewExpressionRequest) *Expression {
	return &Expression{
		NewExpressionRequest: exp,
		Result:               0,
		Status:               "pending",
	}
}

func (e *Expression) IsValid() bool {
	return e.Status != "invalid"
}

// AddExpression выполняет добавление вычисления арифметического выражения
func (s *Service) AddExpression(req *NewExpressionRequest) error {
	exp := NewExpression(req)
	err := s.enqueueExpression(exp)
	if err != nil {
		return err
	}

	s.Logger.Infof("new expression (id: %s): %s", exp.Id, exp.Expression)
	if err := s.generateTasks(exp.Id); err != nil {
		return err
	}

	if s.taskCounter > 0 {
		exp.lastTask = s.tasks[s.taskCounter]
	}

	return nil
}

// completeExpression выполняет всю логику при завершении вычисления выражения
func (s *Service) completeExpression(exp *Expression) {
	exp.Result = exp.lastTask.result
	exp.Status = "done"
	s.clearTasks(exp.lastTask, true)
	exp.lastTask = nil
	s.Logger.Infof("expression (id: %s) done. result: %f", exp.Id, exp.Result)
}

// clearTasks удаляет из памяти задачи, которые требуются для выполнения данной
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

// generateTasks преобразует выражение в ряд задач
func (s *Service) generateTasks(expressionId string) error {
	s.Logger.Debug("generating tasks")

	var cnt uint // подсчёт созданных задач
	exp := s.expressions[expressionId]
	postfix, err := util.ToPostfix(exp.Expression) // получение постфикса (обратная польская запись)

	if err != nil {
		exp.Status = "invalid"
		return err
	}

	//обработка постфикса
	stack := make([]interface{}, 0) // стек для хранения операндов и ссылок на задачи
	for _, token := range postfix {
		if operand, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, operand)
		} else {
			if len(stack) < 2 {
				exp.Status = "invalid"
				return fmt.Errorf("invalid postfix expression")
			}

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
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
					s.Logger.Errorf("expression %v is invalid", exp.Id)
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
		exp.Status = "invalid"
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

// GetExpressionById выполняет получение выражения по Id
func (s *Service) GetExpressionById(id string) (*Expression, bool) {
	exp, exists := s.expressions[id]
	return exp, exists
}

// Task является структурой для задач
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

// GetTask выполняет получение задачи, как правило, самой старой
func (s *Service) GetTask() (*Task, error) {
	task, exists := s.tasks[s.lastTask]
	if !exists {
		return nil, NoTaskError
	}

	exp := s.expressions[task.expressionId]

	// проверка выполнены ли задачи, требуемые для выполнения текущей
	if exp.IsValid() && isTask(task.Arg1) && !task.Arg1.(*Task).isDone {
		return nil, NoTaskError
	}
	if exp.IsValid() && isTask(task.Arg2) && !task.Arg2.(*Task).isDone {
		return nil, NoTaskError
	}

	defer func() { s.lastTask++ }()

	// если выражение невалидное, очистить все задачи
	if !exp.IsValid() {
		s.clearTasks(task, true)
		return s.GetTask()
	}

	// обработка деления на ноль
	val, isFloat := task.Arg2.(float64)
	if (isTask(task.Arg2) && task.Arg2.(*Task).result == 0 || isFloat && val == 0) && task.Operation == "/" {
		exp.Status = "invalid"
		s.Logger.Errorf("expression %v is invalid: division by zero", exp.Id)
		s.clearTasks(task, true)
		return s.GetTask()
	}

	exp.Status = "calculating"

	return task, nil
}

// SetResult выполняет прием результата обработки задачи
func (s *Service) SetResult(id int, result float64) error {
	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("expression %d not found. probably, the expression is invalid", id)
	}
	task.result = result
	task.isDone = true

	exp := s.expressions[task.expressionId]
	s.Logger.Infof("task (id: %d) done. result: %f", id, result)

	// проверка на выполнение всего выражения
	if lastTaskId := exp.lastTask.Id; lastTaskId == task.Id {
		s.completeExpression(exp)
	}
	return nil
}

// GetTaskResponse является основной структурой ответа для получения задачи
type GetTaskResponse struct {
	Task *TaskResponse `json:"task"`
}

// TaskResponse является самим ответом (задачей)
type TaskResponse struct {
	Id            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime uint    `json:"operation_time"`
}

// fillResponse обрабатывает задачу, которую нужно выдать как ответ
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

// GetJSONResponse возвращается json получаемой задачи
func (s *Service) GetJSONResponse(t *Task) ([]byte, error) {
	resp := &GetTaskResponse{
		Task: s.fillResponse(t),
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("не получилось получить json")
	}
	return jsonData, nil
}
