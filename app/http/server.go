package http

import (
	container "clean/app"
	"clean/app/http/middlewares"
	"clean/infra/config"
	"clean/infra/logger"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

func Start() {
	e := echo.New()

	if err := middlewares.Attach(e); err != nil {
		logger.Error("error occur when attaching middlewares", err)
		os.Exit(1)
	}

	container.Init(e.Group("api"))

	port := config.App().Port

	// start http server
	go func() {
		e.Logger.Fatal(e.Start(":" + port))
	}()

	// graceful shutdown
	GracefulShutdown(e)
}

// server will gracefully shutdown within 5 sec
func GracefulShutdown(e *echo.Echo) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = e.Shutdown(ctx)
	logger.Info("server shutdowns gracefully")
}
