package repo

import (
	"orchestrator/internal/models"
)

func (r *UserRepo) CreateUser(user *models.User) error {
	q := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)`

	_, err := r.DB.Exec(q, user.Login, user.Password.Hash)
	return err
}

func (r *UserRepo) GetByLogin(login string) (*models.User, error) {
	q := `
		SELECT id, password FROM users
		WHERE login = $1`

	user := models.User{Password: models.Password{}}

	err := r.DB.QueryRow(q, login).Scan(&user.Id, &user.Password.Hash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
