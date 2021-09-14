package cmd

import (
	"boilerplate/infra/config"
	"boilerplate/infra/conn"
	"boilerplate/infra/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "boilerplate",
		Short: "Pi survey",
		Long:  "A Survey Management system to track survey and survey response",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(seedCmd)
}

// Execute executes the root command
func Execute() {
	config.LoadConfig()
	conn.ConnectDb()
	conn.ConnectMailGun()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Info("about to start the application")
}
