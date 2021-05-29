package domain

type Permission struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"unique" json:"name"`
	Description string `json:"description"`
}
