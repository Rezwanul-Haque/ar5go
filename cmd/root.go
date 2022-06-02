package cmd

import (
	"ar5go/infra/config"
	"ar5go/infra/conn/cache"
	"ar5go/infra/conn/db"
	"ar5go/infra/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "ar5go",
		Short: "implementing clean architecture in golang",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(seedCmd)
}

// Execute executes the root command
func Execute() {
	config.LoadConfig()
	logger.NewLogClient(config.App().LogLevel)
	lc := logger.Client()
	db.NewDbClient(lc)
	cache.NewCacheClient(lc)

	lc.Info("about to start the application")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
