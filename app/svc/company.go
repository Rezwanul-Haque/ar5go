package svc

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infra/errors"
)

type ICompany interface {
	CreateCompanyWithAdminUser(serializers.CompanyReq) (*serializers.CompanyResp, *errors.RestErr)
	GetCompany(companyID uint) (*domain.Company, *errors.RestErr)
}
