package service

import (
	"database/sql"
	"errors"
	"orchestrator/internal/models"
	"orchestrator/pkg/token"
)

func (s *MyService) Register(input UserInput) error {
	newUser := &models.User{
		Login:    input.Login,
		Password: models.Password{},
	}

	err := newUser.Password.Set(input.Password)
	if err != nil {
		return err
	}

	return s.repos.User.CreateUser(newUser)
}

func (s *MyService) Login(input UserInput) (string, error) {
	var password models.Password
	err := password.Set(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repos.User.GetByLogin(input.Login)
	if user == nil || err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows) || user == nil:
			return "", InvalidCreditsError
		default:
			return "", err
		}
	}

	matches, err := user.Password.Matches(input.Password)
	if err != nil {
		return "", err
	}
	if !matches {
		return "", InvalidCreditsError
	}

	return token.New(user.Id, s.Cfg.GetSecret())
}
