package controllers

import (
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/consts"
	"clean/infra/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type location struct {
	hSvc svc.ILocation
}

// NewLocationController will initialize the controllers
func NewLocationController(grp interface{}, ACL func(string) echo.MiddlewareFunc, hSvc svc.ILocation) {
	hc := &location{
		hSvc: hSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/location/history", hc.Create, ACL(consts.PermissionLocationCreate))
}

func (ctr *location) Create(c echo.Context) error {
	var req serializers.LocationHistoryReq

	if err := c.Bind(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		return c.JSON(restErr.Status, restErr)
	}

	saveErr := ctr.hSvc.Create(req)
	if saveErr != nil {
		return c.JSON(saveErr.Status, saveErr)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "location saved"})
}
