package controllers

//import (
//	"clean/app/serializers"
//	"clean/app/svc"
//	"clean/app/utils/errors"
//	"github.com/labstack/echo/v4"
//	"net/http"
//)
//
//type history struct {
//	hSvc svc.IHistory
//}
//
//// NewUsersController will initialize the controllers
//func NewHistoryController(g *echo.Group) {
//	hc := &history{}
//
//	g.POST("/location/history", hc.Create)
//}
//
//func (ctr *history) Create(c echo.Context) error {
//	var req serializers.LocationHistoryReq
//
//	if err := c.Bind(&req); err != nil {
//		restErr := errors.NewBadRequestError("invalid json body")
//		return c.JSON(restErr.Status, restErr)
//	}
//
//	result, saveErr := ctr.hSvc.Create(req)
//	if saveErr != nil {
//		return c.JSON(saveErr.Status, saveErr)
//	}
//
//	return c.JSON(http.StatusCreated, result)
//}
