package repository

import (
	"clean/app/domain"
	"clean/infrastructure/errors"
)

type ICompany interface {
	Save(company *domain.Company) (*domain.Company, *errors.RestErr)
	Get(companyID uint) (*domain.Company, *errors.RestErr)
}
