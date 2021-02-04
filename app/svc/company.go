package svc

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infrastructure/errors"
)

type ICompany interface {
	CreateCompanyWithAdminUser(serializers.CompanyPayload) (*serializers.CompanyResponse, *errors.RestErr)
	GetCompany(companyID uint) (*domain.Company, *errors.RestErr)
}
