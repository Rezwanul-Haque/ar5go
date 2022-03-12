package controllers

import (
	m "ar5go/app/http/middlewares"
	"ar5go/app/serializers"
	"ar5go/app/svc"
	"ar5go/app/utils/consts"
	"ar5go/infra/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type location struct {
	hSvc svc.ILocation
}

// NewLocationController will initialize the controllers
func NewLocationController(grp interface{}, hSvc svc.ILocation) {
	hc := &location{
		hSvc: hSvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/location/history", hc.Create, m.ACL(consts.PermissionLocationCreate))
}

func (ctr location) Create(c echo.Context) error {
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
