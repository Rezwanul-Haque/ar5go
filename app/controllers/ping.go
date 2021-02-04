package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ping struct {}

// NewPingController will initialize the controllers
func NewPingController(g *echo.Group) {
	pc := &ping{}

	g.GET("/ping", pc.Ping)
}

func (clr *ping) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "Pong")
}
