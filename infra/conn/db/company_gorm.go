package db

import (
	"ar5go/app/domain"
	"ar5go/infra/conn/db/models"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
)

func (dc DatabaseClient) SaveCompany(com *domain.Company) (*domain.Company, *errors.RestErr) {
	res := dc.DB.Model(&models.Company{}).Create(&com)

	if res.Error != nil {
		logger.Error("error occurred when create company", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return com, nil
}

func (dc DatabaseClient) GetCompany(companyID uint) (*domain.Company, *errors.RestErr) {
	var com domain.Company
	res := dc.DB.Model(&models.Company{}).Where("id = ?", companyID).First(&com)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("their are no associate company to this user")
	}

	return &com, nil
}
