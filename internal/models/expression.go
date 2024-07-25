package models

// Expression является выражением, которое нужно вычислить
type Expression struct {
	Id         int     `json:"id"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	Status     string  `json:"status"`
	Creator    int64   `json:"-"`
	LastTask   *Task   `json:"-"`
}
