package models

// Expression является выражением, которое нужно вычислить
type Expression struct {
	Id         string  `json:"id"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	Status     string  `json:"status"`
	Creator    int64   `json:"creator"`
	LastTask   *Task   `json:"-"`
}
