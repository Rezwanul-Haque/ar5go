package controllers

import (
	"ar5go/app/domain"
	"ar5go/app/serializers"
	"ar5go/app/svc"
	"ar5go/app/utils/methodsutil"
	"ar5go/infra/errors"
	"fmt"
	"github.com/labstack/echo/v4"
)

func GetUserByAppKey(c echo.Context, uSvc svc.IUsers) (*domain.User, *errors.RestErr) {
	appKey := c.Request().Header.Get("AppKey")

	if methodsutil.IsInvalid(appKey) {
		keyErr := errors.NewBadRequestError(fmt.Sprintf("Appkey: '%s' is missing", appKey))
		return nil, keyErr
	}

	foundUser, getErr := uSvc.GetUserByAppKey(appKey)

	if getErr != nil {
		return nil, getErr
	}

	if appKey != foundUser.AppKey {
		keyErr := errors.NewBadRequestError(fmt.Sprintf("Appkey: '%s' is invalid", appKey))
		return nil, keyErr
	}

	return foundUser, nil
}

func GetUserFromContext(c echo.Context) (*serializers.LoggedInUser, error) {
	user, ok := c.Get("user").(*serializers.LoggedInUser)
	if !ok {
		return nil, errors.ErrNoContextUser
	}

	return user, nil
}
