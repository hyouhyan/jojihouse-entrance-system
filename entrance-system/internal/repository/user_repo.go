package repository

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"jojihouse-entrance-system/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	var totalStayTime_Interval sql.NullString
	user := &model.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(
		&user.ID,
		&user.Name,
		&user.Description,
		&user.Barcode,
		&user.Contact,
		&user.Remaining_entries,
		&user.Registered_at,
		&user.Total_entries,
		&totalStayTime_Interval, //SQLのIntervalで取得
	)
	if err != nil {
		return nil, err
	}

	// Goのtime.Duretionに変換
	user.Total_stay_time = intervalToDuration(totalStayTime_Interval)

	return user, nil
}

func (r *UserRepository) GetUserByBarcode(barcode string) (*model.User, error) {
	var totalStayTime_Interval sql.NullString
	user := &model.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE barcode = $1", barcode).Scan(
		&user.ID,
		&user.Name,
		&user.Description,
		&user.Barcode,
		&user.Contact,
		&user.Remaining_entries,
		&user.Registered_at,
		&user.Total_entries,
		&totalStayTime_Interval, //SQLのIntervalで取得
	)
	if err != nil {
		return nil, err
	}

	// Goのtime.Duretionに変換
	user.Total_stay_time = intervalToDuration(totalStayTime_Interval)

	return user, nil
}

func (r *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	err := r.db.QueryRow(
		"INSERT INTO users (name, description, barcode, contact, remaining_entries, total_entries) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Name,
		user.Description,
		user.Barcode,
		user.Contact,
		user.Remaining_entries,
		user.Total_entries,
	).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET name = $1, description = $2, barcode = $3, contact = $4, remaining_entries = $5, total_entries = $6 WHERE id = $7",
		user.Name,
		user.Description,
		user.Barcode,
		user.Contact,
		user.Remaining_entries,
		user.Total_entries,
		user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// 入場可能回数を減らす
func (r *UserRepository) DecreaseRemainingEntries(id int, count int) (int, int, error) {
	var before int
	var after int

	// remaining_entries を更新しつつ、更新前後の値を取得
	err := r.db.QueryRow(`
		UPDATE users 
		SET remaining_entries = remaining_entries - $1 
		WHERE id = $2
		RETURNING remaining_entries + $1, remaining_entries
	`, count, id).Scan(&before, &after)

	if err != nil {
		return 0, 0, err
	}
	return before, after, nil
}

// 入場可能回数を増やす
func (r *UserRepository) IncreaseRemainingEntries(id int, count int) (int, int, error) {
	var before int
	var after int

	// remaining_entries を更新しつつ、更新前後の値を取得
	err := r.db.QueryRow(`
		UPDATE users
		SET remaining_entries = remaining_entries + $1
		WHERE id = $2
		RETURNING remaining_entries - $1, remaining_entries
	`, count, id).Scan(&before, &after)

	if err != nil {
		return 0, 0, err
	}

	return before, after, nil
}

// 総入場回数を増やす
func (r *UserRepository) IncreaseTotalEntries(id int) error {
	_, err := r.db.Exec("UPDATE users SET total_entries = total_entries + 1 WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// 滞在時間を増やす
func (r *UserRepository) IncreaseTotalStayTime(userID int, stayTime time.Duration) error {
	// SQLのINTERVALに変換
	interval := durationToInterval(stayTime)

	_, err := r.db.Exec("UPDATE users SET total_stay_time = total_stay_time + $1 WHERE id = $2", interval, userID)

	if err != nil {
		return err
	}

	return nil
}

// PostgreSQL INTERVAL → Go time.Duration 変換（分単位）
func intervalToDuration(interval sql.NullString) time.Duration {
	if !interval.Valid {
		return 0
	}
	parsedMinutes, err := strconv.Atoi(strings.Split(interval.String, " ")[0]) // "X minutes" の "X" を取得
	if err != nil {
		return 0
	}
	return time.Duration(parsedMinutes) * time.Minute
}

// Go time.Duration → PostgreSQL INTERVAL 変換（分単位）
func durationToInterval(d time.Duration) string {
	// 分を計算して切り上げ
	minutes := int(math.Ceil(d.Minutes()))

	// SQLのINTERVALに変換
	return fmt.Sprintf("%d minutes", minutes)
}
