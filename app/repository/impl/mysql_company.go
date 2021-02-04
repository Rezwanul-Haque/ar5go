package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/infrastructure/errors"
	"clean/infrastructure/logger"
	"gorm.io/gorm"
)

type company struct {
	*gorm.DB
}

// NewMySqlCompanyRepository will create an object that represent the Company.Repository implementations
func NewMySqlCompanyRepository(db *gorm.DB) repository.ICompany {
	return &company{
		DB: db,
	}
}

func (db *company) Save(com *domain.Company) (*domain.Company, *errors.RestErr) {
	res := db.Model(&domain.Company{}).Create(&com)

	if res.Error != nil {
		logger.Error("error occurred when create company", res.Error)
		return nil, errors.NewInternalServerError("db error")
	}

	return com, nil
}

func (db *company) Get(companyID uint) (*domain.Company, *errors.RestErr) {
	var com domain.Company
	res := db.Model(&domain.Company{}).Where("id = ?", companyID).First(&com)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("their are no associate company to this user")
	}

	return &com, nil
}
