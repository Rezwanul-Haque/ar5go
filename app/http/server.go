package http

import (
	container "ar5go/app"
	"ar5go/app/http/middlewares"
	"ar5go/infra/config"
	"ar5go/infra/logger"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

func Start() {
	e := echo.New()
	lc := logger.Client()

	if err := middlewares.Attach(e, lc); err != nil {
		logger.Client().Error("error occurred when attaching middlewares", err)
		os.Exit(1)
	}

	// routes for documentation
	dg := e.Group("docs")
	dg.GET("/swagger", echo.WrapHandler(middlewares.SwaggerDocs()))
	dg.GET("/redoc", echo.WrapHandler(middlewares.ReDocDocs()))
	dg.GET("/rapidoc", echo.WrapHandler(middlewares.RapiDocs()))
	e.File("/swagger.yaml", "./swagger.yaml")

	// Create a new Prometheus server for metrics using Prometheus Middleware
	echoProm := echo.New()

	middlewares.PrometheusMonitor(echoProm)

	go func() {
		echoProm.Logger.Fatal(echoProm.Start(":" + config.App().MetricsPort))

		// gracefully shutdown metrics server
		GracefulShutdown(echoProm, lc)
	}()

	container.Init(e.Group("api"), lc)

	port := config.App().Port

	// start http server
	go func() {
		e.Logger.Fatal(e.Start(":" + port))
	}()

	// graceful shutdown
	GracefulShutdown(e, lc)
}

// GracefulShutdown server will gracefully shut down within 5 sec
func GracefulShutdown(e *echo.Echo, lc logger.LogClient) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	lc.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = e.Shutdown(ctx)
	lc.Info("server shutdowns gracefully")
}
