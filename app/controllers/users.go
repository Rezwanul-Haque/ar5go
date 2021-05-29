package controllers

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/consts"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type users struct {
	cSvc svc.ICompany
	uSvc svc.IUsers
	lSvc svc.ILocation
}

// NewUsersController will initialize the controllers
func NewUsersController(grp interface{}, ACL func(string) echo.MiddlewareFunc, cSvc svc.ICompany, uSvc svc.IUsers, lSvc svc.ILocation) {
	uc := &users{
		cSvc: cSvc,
		uSvc: uSvc,
		lSvc: lSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/user/signup", uc.Create, ACL(consts.PermissionUserCreate))
	g.GET("/v1/user/resolve", uc.GetAll, ACL(consts.PermissionUserFetchAll))
	g.PATCH("/v1/user", uc.Update, ACL(consts.PermissionUserUpdate))
	g.GET("/v1/user/:user_id/locations", uc.GetUserVisitedLocations, ACL(consts.PermissionUserLocationFetch))
	g.POST("/v1/password/change", uc.ChangePassword)
	g.POST("/v1/password/forgot", uc.ForgotPassword)
	g.POST("/v1/password/verifyreset", uc.VerifyResetPassword)
	g.POST("/v1/password/reset", uc.ResetPassword)
}

func (ctr *users) Create(c echo.Context) error {
	foundUser, getErr := GetUserByAppKey(c, ctr.uSvc)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
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
	var resp serializers.UserResp
	respErr := methodsutil.StructToStruct(result, &resp)
	if respErr != nil {
		return respErr
	}

	return c.JSON(http.StatusCreated, resp)
}

func (ctr *users) GetAll(c echo.Context) error {
	foundUser, getErr := GetUserByAppKey(c, ctr.uSvc)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}
	pagination := GeneratePaginationRequest(c)
	resp, totalPages, getErr := ctr.uSvc.GetUserByCompanyIdAndRole(uint(foundUser.CompanyID), consts.RoleIDSales, pagination)

	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	GeneratePagesPath(c, resp, totalPages)
	return c.JSON(http.StatusOK, &resp)
}

func (ctr *users) Update(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		logger.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}

	var user serializers.UserReq
	if err := c.Bind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	updateErr := ctr.uSvc.UpdateUser(uint(loggedInUser.ID), user)
	if updateErr != nil {
		return c.JSON(updateErr.Status, updateErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityUpdateSuccessMsg("user")})
}

func (ctr *users) GetUserVisitedLocations(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		logger.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}

	pagination := GeneratePaginationRequest(c)
	resp, totalPages, getErr := ctr.lSvc.GetLocationsByUserID(uint(loggedInUser.ID), pagination)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	GeneratePagesPath(c, resp, totalPages)
	return c.JSON(http.StatusOK, &resp)
}

func (ctr *users) ChangePassword(c echo.Context) error {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		logger.Error(err.Error(), err)
		restErr := errors.NewUnauthorizedError("no logged-in user found")
		return c.JSON(restErr.Status, restErr)
	}
	body := &serializers.ChangePasswordReq{}
	if err := c.Bind(&body); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}
	if err = body.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}
	if body.OldPassword == body.NewPassword {
		restErr := errors.NewBadRequestError("password can't be same as old one")
		return c.JSON(restErr.Status, restErr)
	}
	if err := ctr.uSvc.ChangePassword(loggedInUser.ID, body); err != nil {
		switch err {
		case errors.ErrInvalidPassword:
			restErr := errors.NewBadRequestError("old password didn't match")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityChangedSuccessMsg("password")})
}

func (ctr *users) ForgotPassword(c echo.Context) error {
	body := &serializers.ForgotPasswordReq{}

	if err := c.Bind(&body); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := body.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := ctr.uSvc.ForgotPassword(body.Email); err != nil && err == errors.ErrSendingEmail {
		restErr := errors.NewInternalServerError("failed to send password reset email")
		return c.JSON(restErr.Status, restErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Password reset link sent to email"})
}

func (ctr *users) VerifyResetPassword(c echo.Context) error {
	req := &serializers.VerifyResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := req.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := ctr.uSvc.VerifyResetPassword(req); err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidPasswordResetToken:
			restErr := errors.NewUnauthorizedError("failed to send reset_token email")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	return c.JSON(http.StatusOK, "reset token verified")
}

func (ctr *users) ResetPassword(c echo.Context) error {
	req := &serializers.ResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	if err := req.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	verifyReq := &serializers.VerifyResetPasswordReq{
		Token: req.Token,
		ID:    req.ID,
	}

	if err := ctr.uSvc.VerifyResetPassword(verifyReq); err != nil {
		switch err {
		case errors.ErrParseJwt,
			errors.ErrInvalidPasswordResetToken:
			restErr := errors.NewUnauthorizedError("failed to send reset_token email")
			return c.JSON(restErr.Status, restErr)
		default:
			restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
			return c.JSON(restErr.Status, restErr)
		}
	}

	if err := ctr.uSvc.ResetPassword(req); err != nil {
		restErr := errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		return c.JSON(restErr.Status, restErr)
	}

	return c.JSON(http.StatusOK, "password reset successful")
}
