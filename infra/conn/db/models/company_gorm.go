package models

import (
	"gorm.io/gorm"
	"time"
)

type Company struct {
	ID            uint   `gorm:"primarykey"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	Address       string `json:"address"`
	BusinessID    uint   `json:"business_id"`
	NumOfEmployee uint   `json:"num_of_employee"`
	Email         string `json:"email"`
	SnsLink       string `json:"sns_link"`
	Phone         string `json:"phone"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
