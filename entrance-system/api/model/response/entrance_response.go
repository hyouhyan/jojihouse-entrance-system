package response

import (
	"time"
)

type Entrance struct {
	UserID            int       `json:"user_id"`
	UserName          string    `json:"user_name"`
	UserRoles         []string  `json:"user_roles"`
	Time              time.Time `json:"time"`
	AccessType        string    `json:"access_type"`
	Remaining_entries int       `json:"remaining_entries"`
	Number            *int      `json:"number"`
	Total_entries     int       `json:"total_entries"`
}
