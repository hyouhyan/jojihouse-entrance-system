package repository

import (
	"database/sql"
	"jojihouse-entrance-system/internal/model"
	"time"
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

// 在室ユーザー一覧を取得
func (r *UserRepository) GetCurrentUsers() ([]model.CurrentUser, error) {
	var users []model.CurrentUser

	query := `
		SELECT u.id, u.name, c.entered_at
		FROM current_users c
		JOIN users u ON c.user_id = u.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.CurrentUser
		if err := rows.Scan(&user.UserID, &user.Name, &user.EnteredAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// 在室中のユーザーの入室時間を取得
func (r *UserRepository) GetEnteredTime(userID int) (time.Time, error) {
	var enteredAt time.Time

	query := `
		SELECT entered_at
		FROM current_users
		WHERE user_id = $1
	`

	err := r.db.QueryRow(query, userID).Scan(&enteredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return enteredAt, nil
}
