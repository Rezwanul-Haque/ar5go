package domain

import "time"

type LocationHistory struct {
	ID                    uint       `json:"id"`
	CompanyID             uint       `json:"company_id"`
	UserID                uint       `json:"user_id"`
	CheckInTime           time.Time  `json:"check_in_time"`
	CheckOutTime          *time.Time `json:"check_out_time"`
	VisitedCompanyID      uint       `json:"visited_company_id"`
	Latitude              float64    `json:"lat"`
	Longitude             float64    `json:"lon"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`
}
