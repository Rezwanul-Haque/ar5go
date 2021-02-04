package domain

import (
	"time"
)

type Company struct {
	ID            uint       `json:"-"`
	Name          string     `json:"name"`
	Logo          string     `json:"logo"`
	Address       string     `json:"address"`
	BusinessID    uint       `json:"business_id"`
	NumOfEmployee uint       `json:"num_of_employee"`
	Email         string     `json:"email"`
	SnsLink       string     `json:"sns_link"`
	Phone         string     `json:"phone"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
	DeletedAt     *time.Time `json:"-"`
}
