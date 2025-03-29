package repository

import (
	"database/sql"
	"jojihouse-entrance-system/internal/model"
)

type CurrentUsersRepository struct {
	db *sql.DB
}

func NewCurrentUsersRepository(db *sql.DB) *CurrentUsersRepository {
	return &CurrentUsersRepository{db: db}
}

func (r *CurrentUsersRepository) AddUserToCurrentUsers(userInfo model.CurrentUsers) error {
	_, err := r.db.Exec(`
		INSERT INTO current_users (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING;
	`, userInfo.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *CurrentUsersRepository) DeleteUserToCurrentUsers(userInfo model.CurrentUsers) error {
	_, err := r.db.Exec("DELETE FROM current_users WHERE user_id = $1;", userInfo.UserID)
	if err != nil {
		return err
	}

	return nil
}
