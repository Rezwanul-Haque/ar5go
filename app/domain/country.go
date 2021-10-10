package domain

type Country struct {
	ID          uint   `json:"id"`
	Name        string `gorm:"unique" json:"name"`
}
