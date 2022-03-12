package cmd

import (
	"ar5go/infra/conn/db"
	"ar5go/infra/conn/db/seeder"
	"ar5go/infra/logger"

	//"ar5go/infra/conn/db/seeder"
	//"ar5go/infra/logger"
	"fmt"

	"github.com/spf13/cobra"
)

var truncate bool

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Pre populate data",
	Long:  `Pre populate roles, permissions, role_permissions data into mysql`,
	Run:   seed,
}

func seed(cmd *cobra.Command, args []string) {
	// seed roles, permissions, role_permissions
	db.NewDbClient()
	dbc := db.Client()

	truncate, _ = cmd.Flags().GetBool("truncate")
	fmt.Println("truncate=", truncate)

	for _, seed := range seeder.SeedAll() {
		if err := seed.Run(dbc, truncate); err != nil {
			logger.Error(fmt.Sprintf("Running seed '%s', failed with error:", seed.Name), err)
		}
	}
}

func init() {
	seedCmd.PersistentFlags().BoolVarP(&truncate, "truncate", "t", false, "will truncate tables")
}
