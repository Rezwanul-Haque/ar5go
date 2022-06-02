package models

type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"unique" json:"name"`
	DisplayName string `json:"display_name"`
}
