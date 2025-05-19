package repository

import "github.com/jmoiron/sqlx"

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func SetUserPassword()
