package cmd

import (
	server "boilerplate/app/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	// http server setup
	server.Start()
}
