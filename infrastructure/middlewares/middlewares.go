package middlewares

import (
	"clean/infrastructure/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const EchoLogFormat = "time: ${time_rfc3339_nano} || ${method}: ${uri} || status: ${status} || latency: ${latency_human} \n"

// Attach middlewares required for the application, eg: sentry, newrelic etc.
func Attach(e *echo.Echo) error {
	// remove trailing slashes from each requests
	e.Pre(middleware.RemoveTrailingSlash())

	// echo middlewares, todo: add color to the log
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: EchoLogFormat}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(context echo.Context) bool {
			return context.Request().RequestURI == "/metrics"
		},
		Level: 5,
	}))

	e.Use(JWTWithConfig(JWTConfig{
		Skipper: func(context echo.Context) bool {
			switch context.Request().URL.Path {
			case "/api/metrics",
				"/api/v1/login",
				"/api/v1/token/verify",
				"/api/v1/token/refresh",
				"/api/v1/password/forgot",
				"/api/v1/password/verifyreset",
				"/api/v1/password/reset",
				"/api/v1/company/signup":
				return true
			default:
				return false
			}
		},
		SigningKey: []byte(config.Jwt().AccessTokenSecret),
		ContextKey: config.Jwt().ContextKey,
	}))

	return nil
}
