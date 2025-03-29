package repository

import (
	"database/sql"
)

type CurrentUsersRepository struct {
	db *sql.DB
}

func NewCurrentUsersRepository(db *sql.DB) *CurrentUsersRepository {
	return &CurrentUsersRepository{db: db}
}

// 在室ユーザーに追加
func (r *CurrentUsersRepository) AddUserToCurrentUsers(userID int) error {
	_, err := r.db.Exec(`
		INSERT INTO current_users (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING;
	`, userID)
	if err != nil {
		return err
	}

	return nil
}

// 在室ユーザーから削除
func (r *CurrentUsersRepository) DeleteUserToCurrentUsers(userID int) error {
	_, err := r.db.Exec("DELETE FROM current_users WHERE user_id = $1;", userID)
	if err != nil {
		return err
	}

	return nil
}
