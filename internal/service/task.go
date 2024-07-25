package service

import (
	"encoding/json"
	"fmt"
	"orchestrator/internal/models"
	"orchestrator/pkg/util"
	"strconv"
	"sync"
)

func (s *MyService) newTask(arg1, arg2 interface{}, operation string, operationTime uint, expressionId int) *models.Task {
	s.tasks.mu.Lock()
	defer s.tasks.mu.Unlock()

	s.tasks.taskCounter++
	task := &models.Task{
		Id:            s.tasks.taskCounter,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operation,
		OperationTime: operationTime,
		ExpressionId:  expressionId,
		Result:        0,
		IsDone:        false,
		IsCalculating: false,
	}
	s.tasks.tasks[task.Id] = task
	return task
}

func isTask(arg interface{}) bool {
	_, ok := arg.(*models.Task)
	return ok
}

// clearTasks удаляет из памяти задачи, которые требуются для выполнения данной
func (s *MyService) clearTasks(lastTask *models.Task, deleteCurrent bool) {
	if lastTask == nil {
		return
	}
	if isTask(lastTask.Arg1) {
		s.clearTasks(lastTask.Arg1.(*models.Task), true)
		lastTask.Arg1 = nil
	}
	if isTask(lastTask.Arg1) {
		s.clearTasks(lastTask.Arg1.(*models.Task), true)
		lastTask.Arg2 = nil
	}
	if deleteCurrent {
		s.tasks.delete(lastTask.Id)
		s.Logger.Debugf("task %d has been deleted", lastTask.Id)
	}
}

// generateTasks преобразует выражение в ряд задач
func (s *MyService) generateTasks(expressionId int) error {
	s.Logger.Debug("generating tasks")

	exp := s.expressions[expressionId]
	postfix, err := util.ToPostfix(exp.Expression) // получение постфикса (обратная польская запись)

	if err != nil {
		if err := s.invalidate(exp); err != nil {
			return err
		}
		return err
	}

	// Если выражение состоит из одного числа, сразу записываем результат
	if len(postfix) == 1 {
		if operand, err := strconv.ParseFloat(postfix[0], 64); err == nil {
			exp.Result = operand
			exp.Status = "done"
			s.Logger.Infof("no tasks needed for %s. it is done.", exp.Id)
			return nil
		} else {
			exp.Status = "invalid"
			return fmt.Errorf("invalid number format")
		}
	}

	// Обработка постфикса
	var cnt uint                    // подсчёт созданных задач
	stack := make([]interface{}, 0) // стек для хранения операндов и ссылок на задачи
	for _, token := range postfix {
		if operand, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, operand)
		} else {
			if len(stack) < 2 {
				if err := s.invalidate(exp); err != nil {
					return err
				}
				return fmt.Errorf("invalid postfix expression")
			}

			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			cnt++

			var task *models.Task
			switch token {
			case "+":
				task = s.newTask(a, b, "+", s.Cfg.Addition, expressionId)
			case "-":
				task = s.newTask(a, b, "-", s.Cfg.Subtraction, expressionId)
			case "*":
				task = s.newTask(a, b, "*", s.Cfg.Multiplication, expressionId)
			case "/":
				if !isTask(b) && b.(float64) == 0 {
					if err := s.invalidate(exp); err != nil {
						return err
					}
					s.Logger.Errorf("expression %v is invalid", exp.Id)
					return fmt.Errorf("division by zero")
				}
				task = s.newTask(a, b, "/", s.Cfg.Division, expressionId)
			default:
				if err := s.invalidate(exp); err != nil {
					return err
				}
				return fmt.Errorf("invalid operator: %s", token)
			}
			stack = append(stack, task)
		}
	}

	if len(stack) != 1 {
		if err := s.invalidate(exp); err != nil {
			return err
		}
		return fmt.Errorf("invalid postfix expression")
	}

	s.Logger.Infof("successfully created %d tasks", cnt)
	return nil
}

// GetTask выполняет получение задачи, как правило, самой старой
func (s *MyService) GetTask() ([]byte, error) {
	firstLoopFlag := true
	increase := func(flag bool) {
		if firstLoopFlag {
			s.tasks.mu.Lock()
			defer s.tasks.mu.Unlock()
			s.tasks.lastTask++
		}
	}

	for i := s.tasks.lastTask; i <= s.tasks.taskCounter; i++ {
		task, exists := s.tasks.get(i)

		if !exists || task.IsDone || task.IsCalculating {
			firstLoopFlag = false
			continue
		}

		exp := s.expressions[task.ExpressionId]

		// проверка выполнены ли задачи, требуемые для выполнения текущей
		if isValid(exp) && isTask(task.Arg1) && !task.Arg1.(*models.Task).IsDone || isValid(exp) && isTask(task.Arg2) && !task.Arg2.(*models.Task).IsDone {
			firstLoopFlag = false
			continue
		}

		// если выражение невалидное, очистить все задачи
		if !isValid(exp) {
			s.clearTasks(task, true)
			increase(firstLoopFlag)
			return s.GetTask()
		}

		// обработка деления на ноль
		val, isFloat := task.Arg2.(float64)
		if (isTask(task.Arg2) && task.Arg2.(*models.Task).Result == 0 || isFloat && val == 0) && task.Operation == "/" {
			if err := s.invalidate(exp); err != nil {
				return nil, err
			}
			s.Logger.Errorf("expression %v is invalid: division by zero", exp.Id)
			s.clearTasks(task, true)
			increase(firstLoopFlag)
			return s.GetTask()
		}

		exp.Status = "calculating"
		task.IsCalculating = true
		increase(firstLoopFlag)

		return s.getJSONResponse(task)
	}

	return nil, NoTaskError
}

// SetTaskResult выполняет прием результата обработки задачи
func (s *MyService) SetTaskResult(id int, result float64) error {
	task, exists := s.tasks.get(id)
	if !exists {
		return fmt.Errorf("expression %d not found. probably, the expression is invalid", id)
	}
	task.Result = result
	task.IsDone = true

	exp := s.expressions[task.ExpressionId]
	s.Logger.Infof("task (id: %d) done. Result: %f", id, result)

	// проверка на выполнение всего выражения
	if lastTaskId := exp.LastTask.Id; lastTaskId == task.Id {
		return s.completeExpression(exp)
	}
	return nil
}

type taskQueue struct {
	tasks       map[int]*models.Task
	taskCounter int // Переменная, для хранения id каждого новой задачи
	lastTask    int // Переменная, для хранения id последней выполненной задачи
	mu          *sync.RWMutex
}

func newTaskQueue() *taskQueue {
	return &taskQueue{
		tasks:       make(map[int]*models.Task),
		taskCounter: 0,
		lastTask:    0,
		mu:          &sync.RWMutex{},
	}
}

func (q *taskQueue) get(id int) (*models.Task, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	task, ok := q.tasks[id]
	return task, ok
}

func (q *taskQueue) delete(id int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.tasks, id)
}

// fillResponse обрабатывает задачу, которую нужно выдать как ответ
func (s *MyService) fillResponse(task *models.Task) *TaskResponse {
	var arg1, arg2 float64

	if isTask(task.Arg1) {
		arg1 = task.Arg1.(*models.Task).Result
	} else {
		arg1 = task.Arg1.(float64)
	}

	if isTask(task.Arg2) {
		arg2 = task.Arg2.(*models.Task).Result
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

// getJSONResponse возвращается json получаемой задачи
func (s *MyService) getJSONResponse(t *models.Task) ([]byte, error) {
	resp := &GetTaskResponse{
		Task: s.fillResponse(t),
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("invalid json")
	}

	return jsonData, nil
}

func (s *MyService) LoadTasks() error {
	for _, exp := range s.expressions {
		if exp.Status == "pending" {
			if err := s.generateTasks(exp.Id); err != nil {
				return err
			}
			if s.tasks.taskCounter > 0 {
				exp.LastTask, _ = s.tasks.get(s.tasks.taskCounter)
			}
		}
	}

	return nil
}
