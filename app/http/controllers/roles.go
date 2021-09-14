package controllers

import (
	"boilerplate/app/serializers"
	"boilerplate/app/svc"
	"boilerplate/app/utils/consts"
	"boilerplate/app/utils/msgutil"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type roles struct {
	rsvc svc.IRoles
}

// NewRolesController will initialize the controllers
func NewRolesController(grp interface{}, ACL func(string) echo.MiddlewareFunc, rsvc svc.IRoles) {
	rc := &roles{
		rsvc: rsvc,
	}

	g := grp.(*echo.Group)

	g.POST("/v1/role", rc.CreateRole, ACL(consts.PermissionRoleCrud))
	g.PATCH("/v1/role/:role_id", rc.UpdateRole, ACL(consts.PermissionRoleCrud))
	g.DELETE("/v1/role/:role_id", rc.DeleteRole, ACL(consts.PermissionRoleCrud))
	g.GET("/v1/role", rc.ListRoles, ACL(consts.PermissionRoleFetchAll))

	g.POST("/v1/role/:role_id/permissions", rc.SetRolePermissions, ACL(consts.PermissionRoleCrud))
	g.GET("/v1/role/:role_id/permissions", rc.GetRolePermissions, ACL(consts.PermissionRoleCrud))
}

func (ctr *roles) CreateRole(c echo.Context) error {
	var roleToCreate *serializers.RoleReq
	var err error

	if err = c.Bind(&roleToCreate); err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("bind role body to struct"), err)
		restErr := errors.NewBadRequestError(errors.ErrCheckParamBodyHeader)
		return c.JSON(restErr.Status, restErr)
	}

	logger.Info(fmt.Sprintf("%+v", roleToCreate))
	if err = roleToCreate.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	resp, createErr := ctr.rsvc.CreateRole(roleToCreate)
	if createErr != nil {
		return c.JSON(createErr.Status, createErr)
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctr *roles) UpdateRole(c echo.Context) error {
	var roleToUpdate serializers.RoleReq

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if err = c.Bind(&roleToUpdate); err != nil {
		logger.Error(msgutil.EntityBindToStructFailedMsg("update role"), err)
		restErr := errors.NewBadRequestError(errors.ErrCheckParamBodyHeader)
		return c.JSON(restErr.Status, restErr)
	}

	logger.Info(fmt.Sprintf("%+v", roleToUpdate))
	if err = roleToUpdate.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if upErr := ctr.rsvc.UpdateRole(uint(roleID), roleToUpdate); upErr != nil {
		return c.JSON(upErr.Status, upErr)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityUpdateSuccessMsg("role")})
}

func (ctr *roles) DeleteRole(c echo.Context) error {
	id := c.Param("role_id")

	valErr := v.Validate(id, v.Required, v.Match(regexp.MustCompile("^[0-9]+$")).Error("invalid role id"))
	if valErr != nil {
		restErr := errors.NewBadRequestError(valErr.Error())
		return c.JSON(restErr.Status, restErr)
	}

	roleID, _ := strconv.Atoi(id)

	if delErr := ctr.rsvc.DeleteRole(uint(roleID)); delErr != nil {
		return c.JSON(delErr.Status, delErr)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": msgutil.EntityDeleteSuccessMsg("role")})
}

func (ctr *roles) ListRoles(c echo.Context) error {
	res, err := ctr.rsvc.ListRoles()

	if err != nil {
		return c.JSON(err.Status, err)
	}
	return c.JSON(http.StatusOK, res)
}

func (ctr *roles) SetRolePermissions(c echo.Context) error {
	var rp serializers.RolePermissionsReq
	var err error

	if err = c.Bind(&rp); err != nil {
		logger.Error(msgutil.EntityBindToStructFailedMsg("role & permission"), err)
		restErr := errors.NewBadRequestError(errors.ErrCheckParamBodyHeader)
		return c.JSON(restErr.Status, restErr)
	}

	roleID, _ := strconv.Atoi(c.Param("role_id"))
	rp.RoleID = roleID

	logger.Info(fmt.Sprintf("%+v", rp))
	
	if err = rp.Validate(); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		return c.JSON(restErr.Status, restErr)
	}

	if setErr := ctr.rsvc.SetRolePermissions(&rp); setErr != nil {
		return c.JSON(setErr.Status, msgutil.EntityCreationFailedMsg("roles permission"))
	}
	return c.JSON(http.StatusOK, rp)
}

func (ctr *roles) GetRolePermissions(c echo.Context) error {
	id := c.Param("role_id")

	valErr := v.Validate(id, v.Required, v.Match(regexp.MustCompile("^[0-9]+$")).Error("invalid role id"))
	if valErr != nil {
		restErr := errors.NewBadRequestError(valErr.Error())
		return c.JSON(restErr.Status, restErr)
	}

	roleID, _ := strconv.Atoi(id)

	res, getErr := ctr.rsvc.GetRolePermissions(roleID)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}
	return c.JSON(http.StatusOK, res)
}
