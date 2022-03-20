package controllers

import (
	"ar5go/app/svc"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type system struct {
	lc  logger.LogClient
	svc svc.ISystem
}

// NewSystemController will initialize the controllers
func NewSystemController(grp interface{}, lc logger.LogClient, sysSvc svc.ISystem) {
	pc := &system{
		lc:  lc,
		svc: sysSvc,
	}

	g := grp.(*echo.Group)

	g.GET("/v1", pc.Root)
	g.GET("/v1/h34l7h", pc.Health)
}

// Root will let you see what you can slash 🐲
func (ctr *system) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ar5go architecture backend! let's play!!"})
}

// swagger:route GET /v1/h34l7h Health will let you know the heart beats ❤️
// Return a message
// responses:
//	200: genericSuccessResponse

// Health will let you know the heart beats ❤️
func (ctr *system) Health(c echo.Context) error {
	resp, err := ctr.svc.GetHealth()
	if err != nil {
		ctr.lc.Error(fmt.Sprintf("%+v", resp), err)
		return c.JSON(http.StatusInternalServerError, errors.ErrSomethingWentWrong)
	}
	return c.JSON(http.StatusOK, resp)
}
