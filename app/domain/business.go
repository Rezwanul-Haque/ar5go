package domain

type Business struct {
	ID          uint   `json:"id"`
	Name        string `gorm:"unique" json:"name"`
	Description string `json:"description"`
}
