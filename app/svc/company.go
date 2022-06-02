package svc

import (
	"ar5go/app/domain"
	"ar5go/app/serializers"
	"ar5go/infra/errors"
)

type ICompany interface {
	CreateCompanyWithAdminUser(serializers.CompanyReq) (*serializers.CompanyResp, *errors.RestErr)
	GetCompany(companyID uint) (*domain.Company, *errors.RestErr)
}
