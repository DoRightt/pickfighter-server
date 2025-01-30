package cmd

import (
	"context"
	"fmt"

	"github.com/DoRightt/pickfighter-server/fighters/internal/repository/psql"
	migrations "github.com/DoRightt/pickfighter-server/fighters/migrations/init"
	"github.com/DoRightt/pickfighter-server/fighters/pkg/cfg"
	logs "github.com/DoRightt/pickfighter-server/pkg/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(initSchemaCmd)
}

var repoCmd = &cobra.Command{
	Use:          "repo",
	Short:        "helps to communicate with application PostgreSQL database",
	Long:         ``,
	SilenceUsage: true,
}

var initSchemaCmd = &cobra.Command{
	Use:   "init-db",
	Short: "DB initialization",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		config := cfg.ViperTestPostgres()
		db, err := psql.New(ctx, config)
		if err != nil {
			logs.Fatalf("Unable to connect PostgreSQL: %s", err)
		}
		defer db.GracefulShutdown()

		err = migrations.InitFightersSchema(ctx, db)
		if err != nil {
			logs.Fatalf("Error while initializing database : %v", err)
		}
		fmt.Println("Database successfully initalized")
	},
}
