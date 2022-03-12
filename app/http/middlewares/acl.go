package middlewares

import (
	"ar5go/app/serializers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ACL(permissionToCheck string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*serializers.LoggedInUser)
			if !ok {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "no logged-in user found"})
			}

			if user.HasPermission(permissionToCheck) {
				return next(c)
			}

			return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "access forbidden"})
		}
	}
}
