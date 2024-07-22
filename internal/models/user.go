package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64    `json:"id"`
	Login    string   `json:"login"`
	Password Password `json:"-"`
}

type Password struct {
	plain string
	Hash  []byte
}

func (p *Password) Set(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return err
	}
	p.plain = plain
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plain))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
