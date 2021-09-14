package controllers

import (
	"boilerplate/app/domain"
	"boilerplate/app/serializers"
	"boilerplate/app/svc"
	"boilerplate/app/utils/consts"
	"boilerplate/app/utils/methodutil"
	"boilerplate/app/utils/msgutil"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

type users struct {
	uSvc svc.IUsers
}

// NewUsersController will initialize the controllers
func NewUsersController(grp interface{}, ACL func(string) echo.MiddlewareFunc, uSvc svc.IUsers) {
	uc := &users{
		uSvc: uSvc,
	}

	g := grp.(*echo.Group)

	g.GET("/v1/check/valid/username", uc.UserNameIsUnique)
	g.GET("/v1/check/valid/email", uc.EmailIsUnique)
	g.POST("/v1/user/signup", uc.Create)
	g.PATCH("/v1/user", uc.Update)
	g.POST("/v1/password/change", uc.ChangePassword)
	g.POST("/v1/password/forgot", uc.ForgotPassword)
	g.POST("/v1/password/verifyreset", uc.VerifyResetPassword)
	g.POST("/v1/password/reset", uc.ResetPassword)
}

func (ctr *users) UserNameIsUnique(c echo.Context) error {

	req := c.QueryParam("user_name")
	if req == "" {
		return c.JSON(http.StatusBadRequest, "requested username should not be empty")
	}

	logger.InfoAsJson("username payload", req)

	if err := ctr.uSvc.UserNameIsUnique(req); err != nil {
		logger.ErrorAsJson("username uniqueness", err)
		restErr := errors.NewAlreadyExistError(errors.ErrUserNameNotUnique)
		return c.JSON(restErr.Status, restErr)
	}
	return c.JSON(http.StatusOK, "Requested username is available")
}

func (ctr *users) EmailIsUnique(c echo.Context) error {
	req := &serializers.EmailIsUnique{}
	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	logger.InfoAsJson("email uniquness payload", req)
	if payloadErr := req.Validate(); payloadErr != nil {
		return payloadErr
	}
	if err := ctr.uSvc.EmailIsUnique(req); err != nil {
		logger.ErrorAsJson("email uniqueness", err)
		restErr := errors.NewAlreadyExistError(errors.ErrEmailIsUnique)
		return c.JSON(restErr.Status, restErr)
	}
	return c.JSON(http.StatusOK, "Requested email is available")
}

func (ctr *users) Create(c echo.Context) error {

	var user domain.User

	if err := c.Bind(&user); err != nil {
		logger.Error("failed to parse request body", err)
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}
	user.RoleID = consts.RoleIDAdmin

	if payloadErr := user.Validate(); payloadErr != nil {
		logger.ErrorAsJson("failed to validate request body", payloadErr)
		restErr := errors.NewBadRequestError(errors.ErrRecordNotvalid)
		return c.JSON(restErr.Status, restErr)
	}

	logger.InfoAsJson("user payload", user)
	result, saveErr := ctr.uSvc.CreateUser(user)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}
	var resp serializers.UserResp
	respErr := methodutil.StructToStruct(result, &resp)
	if respErr != nil {
		return respErr
	}
	return c.JSON(http.StatusCreated, resp)
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

	logger.InfoAsJson("user update payload", user)
	updateErr := ctr.uSvc.UpdateUser(uint(loggedInUser.ID), user)
	if updateErr != nil {
		return c.JSON(updateErr.Status, updateErr)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityUpdateSuccessMsg("user")})
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

	logger.InfoAsJson("change password payload", body)
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

	logger.InfoAsJson("forgot password payload", body)

	if err := body.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err := ctr.uSvc.ForgotPassword(body.Email); err != nil {
		if err == errors.ErrSendingEmail {
			logger.Error(msgutil.EntityGenericFailedMsg("failed to send email"), err)
		}
		logger.Error(msgutil.EntityGenericFailedMsg("failed to send user signup email"), err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "password reset link sent to email"})
}

func (ctr *users) VerifyResetPassword(c echo.Context) error {
	req := &serializers.VerifyResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	logger.InfoAsJson("verify password payload", req)
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
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "reset token verified"})
}

func (ctr *users) ResetPassword(c echo.Context) error {
	req := &serializers.ResetPasswordReq{}

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	logger.InfoAsJson("reset password payload", req)
	if err := req.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	verifyReq := &serializers.VerifyResetPasswordReq{
		Token: req.Token,
		ID:    req.ID,
	}
	logger.InfoAsJson("verify request for reset password", verifyReq)

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
