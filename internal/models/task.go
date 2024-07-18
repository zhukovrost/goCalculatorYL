package models

// Task является структурой для задач
type Task struct {
	Id            int
	Arg1          interface{}
	Arg2          interface{}
	Operation     string
	OperationTime uint
	Result        float64
	ExpressionId  string
	IsDone        bool
	IsCalculating bool
}
