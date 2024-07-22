package repo

import "database/sql"

type Repos struct {
	User       UserRepo
	Expression ExpressionRepo
}

func NewRepos(db *sql.DB) *Repos {
	return &Repos{
		User:       UserRepo{db},
		Expression: ExpressionRepo{db},
	}
}

type UserRepo struct {
	DB *sql.DB
}

type ExpressionRepo struct {
	DB *sql.DB
}
