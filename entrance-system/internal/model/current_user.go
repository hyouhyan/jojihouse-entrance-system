package model

import "time"

type CurrentUser struct {
	UserID    int
	Name      string
	EnteredAt time.Time
}
