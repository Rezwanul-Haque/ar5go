package controllers

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/consts"
	"clean/infrastructure/errors"
	"clean/infrastructure/methodsutil"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type users struct {
	cSvc svc.ICompany
	uSvc svc.IUsers
}

// NewUsersController will initialize the controllers
func NewUsersController(g *echo.Group, cSvc svc.ICompany, uSvc svc.IUsers) {
	uc := &users{
		cSvc: cSvc,
		uSvc: uSvc,
	}

	g.POST("/user/signup", uc.Create)
	g.GET("/user/resolve", uc.GetAll)
}

func (ctr *users) Create(c echo.Context) error {
	appKey := c.Request().Header.Get("AppKey")
	if methodsutil.IsInvalid(appKey) {
		keyErr := errors.NewInternalServerError(fmt.Sprintf("Appkey: '%s' is missing", appKey))
		return c.JSON(keyErr.Status, keyErr)
	}

	foundUser, getErr := ctr.uSvc.GetUserByAppKey(appKey)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	if appKey != foundUser.AppKey {
		keyErr := errors.NewInternalServerError(fmt.Sprintf("Appkey: '%s' is invalid", appKey))
		return c.JSON(keyErr.Status, keyErr)
	}

	var user domain.User

	if err := c.Bind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), 8)
	*user.Password = string(hashedPass)
	user.CompanyID = foundUser.CompanyID
	user.RoleID = consts.RoleIDSales

	result, saveErr := ctr.uSvc.CreateUser(user)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}
	return c.JSON(http.StatusCreated, result)
}

func (ctr *users) GetAll(c echo.Context) error {
	appKeyHeader := c.Request().Header.Get("AppKey")
	if methodsutil.IsInvalid(appKeyHeader) {
		keyErr := errors.NewInternalServerError(fmt.Sprintf("Appkey: '%s' is missing", appKeyHeader))
		return c.JSON(keyErr.Status, keyErr)
	}

	foundUser, getErr := ctr.uSvc.GetUserByAppKey(appKeyHeader)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	if appKeyHeader != foundUser.AppKey {
		keyErr := errors.NewInternalServerError(fmt.Sprintf("Appkey: '%s' is invalid", appKeyHeader))
		return c.JSON(keyErr.Status, keyErr)
	}

	var result serializers.ResolveResponse

	company, getErr := ctr.cSvc.GetCompany(foundUser.CompanyID)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	result.CompanyName = company.Name
	result.CompanyID = foundUser.CompanyID

	subordinates, getErr := ctr.uSvc.GetUserByCompanyIdAndRole(foundUser.CompanyID, consts.RoleIDSales)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	result.Subordinates = subordinates

	return c.JSON(http.StatusOK, result)
}

func GetUserFromContext(c echo.Context) (*serializers.LoggedInUser, error) {
	user, ok := c.Get("user").(*serializers.LoggedInUser)
	if !ok {
		return nil, errors.ErrNoContextUser
	}

	return user, nil
}
