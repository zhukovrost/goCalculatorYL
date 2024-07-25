package service

import (
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"orchestrator/internal/config"
	"orchestrator/internal/models"
	"orchestrator/internal/repo"
)

var (
	NoTaskError         = errors.New("no task")
	InvalidCreditsError = errors.New("invalid credits")
)

/*
type Service interface {
	GetTask() ([]byte, error)
	SetTaskResult(id int, result float64) error
	GetExpressions() []*Expression
	GetExpressionById(id string) (*Expression, bool)
	AddExpression(req *NewExpressionRequest) error
}
*/

// ===== MyService block =====

type MyService struct {
	Cfg         *config.Config
	DB          *sql.DB
	Logger      *logrus.Logger
	repos       *repo.Repos
	expressions map[int]*models.Expression
	tasks       *taskQueue
	LastId      int
}

func New(cfg *config.Config, db *sql.DB, logger *logrus.Logger) (*MyService, error) {
	repos := repo.NewRepos(db)
	srv := &MyService{
		Cfg:         cfg,
		Logger:      logger,
		DB:          db,
		expressions: make(map[int]*models.Expression),
		repos:       repos,
		tasks:       newTaskQueue(),
	}

	exps, last, err := repos.Expression.GetAll()

	if err != nil {
		return srv, err
	}

	srv.expressions = exps
	srv.LastId = last

	return srv, nil
}

// NewExpressionRequest является входными данными при приёме нового выражения
type NewExpressionRequest struct {
	Expression string `json:"expression"`
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

// CalculationResult является структурой получения результата вычисления задачи
type CalculationResult struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

type UserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
