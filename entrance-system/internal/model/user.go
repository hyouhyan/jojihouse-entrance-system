package model

import "time"

type User struct {
	ID                int
	Name              string
	Description       string
	Barcode           string
	Contact           string
	Remaining_entries int
	Registered_at     time.Time
	Total_entries     int
	Total_stay_time   time.Duration
}
