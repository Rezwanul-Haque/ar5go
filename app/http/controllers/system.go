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
	svc svc.ISystem
}

// NewSystemController will initialize the controllers
func NewSystemController(grp interface{}, sysSvc svc.ISystem) {
	pc := &system{
		svc: sysSvc,
	}

	g := grp.(*echo.Group)

	g.GET("/v1", pc.Root)
	g.GET("/v1/h34l7h", pc.Health)
}

// Root will let you see what you can slash 🐲
func (sh *system) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ar5go architecture backend! let's play!!"})
}

// Health will let you know the heart beats ❤️
func (sys *system) Health(c echo.Context) error {
	resp, err := sys.svc.GetHealth()
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", resp), err)
		return c.JSON(http.StatusInternalServerError, errors.ErrSomethingWentWrong)
	}
	return c.JSON(http.StatusOK, resp)
}
