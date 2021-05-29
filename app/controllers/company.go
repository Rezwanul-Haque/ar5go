package controllers

import (
	"clean/app/serializers"
	"clean/app/svc"
	"clean/infra/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type company struct {
	cSvc svc.ICompany
}

// NewCompanyController will initialize the controllers
func NewCompanyController(grp interface{}, cSvc svc.ICompany) {
	cc := &company{
		cSvc: cSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/company/signup", cc.CreateWithAdminUser)
}

func (ctr *company) CreateWithAdminUser(c echo.Context) error {
	var cp serializers.CompanyReq

	if err := c.Bind(&cp); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	result, saveErr := ctr.cSvc.CreateCompanyWithAdminUser(cp)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusCreated, result)
}
