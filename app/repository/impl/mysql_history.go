package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/infrastructure/errors"
	"clean/infrastructure/logger"
	"gorm.io/gorm"
)

type history struct {
	*gorm.DB
}

// NewMySqlHistoryRepository will create an object that represent the Company.Repository implementations
func NewMySqlHistoryRepository(db *gorm.DB) repository.IHistory {
	return &history{
		DB: db,
	}
}

func (r *history) Save(history *domain.LocationHistory) *errors.RestErr {
	res := r.DB.Model(&domain.LocationHistory{}).Create(&history)

	if res.Error != nil {
		logger.Error("error occurred when saving location history", res.Error)
		return errors.NewInternalServerError("db error")
	}

	return nil
}
