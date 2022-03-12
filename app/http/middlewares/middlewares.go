package middlewares

import (
	"ar5go/infra/config"
	"net/http"

	openMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo-contrib/prometheus"
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
				"/api/v1",
				"/api/v1/h34l7h",
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

// PrometheusMonitor will start a middleware which will be
// exposed /metrics handler to be used by prometheus
func PrometheusMonitor(e *echo.Echo) {
	e.HideBanner = true
	prom := prometheus.NewPrometheus("ar5go", nil)

	// Scrape metrics from Main Server
	e.Use(prom.HandlerFunc)
	// Setup metrics endpoint at application server
	prom.SetMetricsPath(e)
}

func SwaggerDocs() http.Handler {
	opts := openMiddleware.SwaggerUIOpts{
		Path:    "docs/swagger",
		SpecURL: "/swagger.yaml",
	}
	return openMiddleware.SwaggerUI(opts, nil)
}

func ReDocDocs() http.Handler {
	opts := openMiddleware.RedocOpts{
		Path:    "docs/redoc",
		SpecURL: "/swagger.yaml",
	}
	return openMiddleware.Redoc(opts, nil)
}

func RapiDocs() http.Handler {
	opts := openMiddleware.RapiDocOpts{
		Path:    "docs/rapidoc",
		SpecURL: "/swagger.yaml",
	}
	return openMiddleware.RapiDoc(opts, nil)
}
