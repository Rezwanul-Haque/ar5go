package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/app/serializers"
	"ar5go/app/svc"
	"ar5go/app/utils/consts"
	"ar5go/app/utils/methodsutil"
	"ar5go/app/utils/msgutil"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type company struct {
	lc    logger.LogClient
	crepo repository.ICompany
	urepo repository.IUsers
}

func NewCompanyService(lc logger.LogClient, crepo repository.ICompany, urepo repository.IUsers) svc.ICompany {
	return &company{
		lc:    lc,
		crepo: crepo,
		urepo: urepo,
	}
}

func (c *company) CreateCompanyWithAdminUser(cp serializers.CompanyReq) (*serializers.CompanyResp, *errors.RestErr) {
	cp.TrimRequestBody()

	var companyObj domain.Company
	if err := methodsutil.StructToStruct(cp, &companyObj); err != nil {
		c.lc.Error(msgutil.EntityStructToStructFailedMsg("admin user"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	companyResult, createErr := c.crepo.SaveCompany(&companyObj)
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

	userResult, createErr := c.urepo.SaveUser(&user)
	if createErr != nil {
		return nil, createErr
	}

	var resp serializers.CompanyResp

	if err := methodsutil.StructToStruct(companyResult, &resp); err != nil {
		c.lc.Error(msgutil.EntityStructToStructFailedMsg("company"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	if err := methodsutil.StructToStruct(userResult, &resp.Admin); err != nil {
		c.lc.Error(msgutil.EntityStructToStructFailedMsg("admin user"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	resp.Admin.Password = nil

	return &resp, nil
}

func (c *company) GetCompany(companyID uint) (*domain.Company, *errors.RestErr) {
	resp, getErr := c.crepo.GetCompany(companyID)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}
