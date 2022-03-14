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

	if err := middlewares.Attach(e); err != nil {
		logger.Client().Error("error occur when attaching middlewares", err)
		os.Exit(1)
	}

	// Create Prometheus server and Middleware
	echoProm := echo.New()

	middlewares.PrometheusMonitor(echoProm)

	go func() {
		echoProm.Logger.Fatal(echoProm.Start(":" + config.App().MetricsPort))

		// gracefully shutdown metrics server
		GracefulShutdown(echoProm)
	}()

	container.Init(e.Group("api"), logger.Client())

	port := config.App().Port

	// start http server
	go func() {
		e.Logger.Fatal(e.Start(":" + port))
	}()

	// graceful shutdown
	GracefulShutdown(e)
}

// GracefulShutdown server will gracefully shut down within 5 sec
func GracefulShutdown(e *echo.Echo) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	logger.Client().Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = e.Shutdown(ctx)
	logger.Client().Info("server shutdowns gracefully")
}
