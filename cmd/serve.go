package cmd

import (
	container "clean/app"
	"clean/app/controllers/middlewares"
	"clean/infra/config"
	"clean/infra/logger"
	"context"

	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	// http server setup
	e := echo.New()

	if err := middlewares.Attach(e); err != nil {
		logger.Error("error occur when attaching middlewares", err)
		os.Exit(1)
	}

	g := e.Group("api")

	container.Init(g)

	port := config.App().Port

	// start http server
	go func() {
		e.Logger.Fatal(e.Start(":" + port))
	}()

	// graceful shutdown setup
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = e.Shutdown(ctx)
	logger.Info("server shutdowns gracefully")
}
