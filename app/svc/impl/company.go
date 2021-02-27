package svc

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/consts"
	"clean/infra/errors"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type company struct {
	crepo repository.ICompany
	urepo repository.IUsers
}

func NewCompanyService(crepo repository.ICompany, urepo repository.IUsers) svc.ICompany {
	return &company{
		crepo: crepo,
		urepo: urepo,
	}
}

func (c *company) CreateCompanyWithAdminUser(cp serializers.CompanyReq) (*serializers.CompanyResp, *errors.RestErr) {
	cp.TrimRequestBody()

	var companyObj domain.Company

	jsonData, _ := json.Marshal(cp)
	_ = json.Unmarshal(jsonData, &companyObj)

	companyResult, createErr := c.crepo.Save(&companyObj)
	if createErr != nil {
		return nil, createErr
	}

	var user domain.User
	user.CompanyID = companyResult.ID
	user.FirstName = "admin"
	user.LastName = "user"
	user.Password = cp.Password
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), 8)
	*user.Password = string(hashedPass)
	user.AppKey = uuid.New().String()
	user.RoleID = consts.RoleIDAdmin
	user.Phone = companyResult.Phone
	user.Email = companyResult.Email

	userResult, createErr := c.urepo.Save(&user)
	if createErr != nil {
		return nil, createErr
	}

	var resp serializers.CompanyResp
	jsonData, _ = json.Marshal(companyResult)
	_ = json.Unmarshal(jsonData, &resp)
	resp.Admin = *userResult

	return &resp, nil
}

func (c *company) GetCompany(companyID uint) (*domain.Company, *errors.RestErr) {
	resp, getErr := c.crepo.Get(companyID)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}
