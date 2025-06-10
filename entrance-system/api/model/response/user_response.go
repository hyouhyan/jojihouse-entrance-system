package response

import (
	"time"
)

type User struct {
	ID                *int       `json:"id"`
	Name              *string    `json:"name"`
	Description       *string    `json:"description"`
	Barcode           *string    `json:"barcode"`
	Contact           *string    `json:"contact"`
	Remaining_entries *int       `json:"remaining_entries"`
	Registered_at     *time.Time `json:"registered_at"`
	Total_entries     *int       `json:"total_entries"`
	Allergy           *string    `json:"allergy"`
	Number            *int       `json:"number"`
	Roles             []Role     `json:"roles"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
