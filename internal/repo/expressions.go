package repo

import (
	"database/sql"
	"errors"
	"orchestrator/internal/models"
)

func (r *ExpressionRepo) GetAll() (map[int]*models.Expression, int, error) {
	q := `SELECT id, expression, creator, result, status FROM expressions`
	res := make(map[int]*models.Expression)
	var last int

	rows, err := r.DB.Query(q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		exp := models.Expression{}
		err = rows.Scan(&exp.Id, &exp.Expression, &exp.Creator, &exp.Result, &exp.Status)
		if err != nil {
			return nil, 0, err
		}
		if exp.Id > last {
			last = exp.Id
		}
		res[exp.Id] = &exp
	}

	return res, last, nil
}

func (r *ExpressionRepo) Add(e *models.Expression) error {
	q := `INSERT INTO expressions(id, expression, creator, result, status) 
	VALUES ($1, $2, $3, $4, $5)`

	_, err := r.DB.Exec(q, e.Id, e.Expression, e.Creator, e.Result, e.Status)

	return err
}

func (r *ExpressionRepo) Update(e *models.Expression) error {
	q := `UPDATE expressions 
	SET result = $1, status = $2
	WHERE id = $3`

	_, err := r.DB.Exec(q, e.Result, e.Status, e.Id)
	return err
}
