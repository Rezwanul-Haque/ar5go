package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ping struct{}

// NewPingController will initialize the controllers
func NewPingController(grp interface{}) {
	pc := &ping{}

	g := grp.(*echo.Group)

	g.GET("/v1/ping", pc.Ping)
}

func (clr *ping) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "pong"})
}
