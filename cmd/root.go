package cmd

import (
	"clean/infrastructure/config"
	"clean/infrastructure/conn"
	"clean/infrastructure/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "clean code",
		Short: "implementing clean architecture in golang",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// Execute executes the root command
func Execute() {
	config.LoadConfig()
	conn.ConnectDb()
	conn.ConnectRedis()

	logger.Info("about to start the application")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
