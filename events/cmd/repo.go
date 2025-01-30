package cmd

import (
	"context"
	"fmt"

	"github.com/DoRightt/pickfighter-server/events/internal/repository/psql"
	migrations "github.com/DoRightt/pickfighter-server/events/migrations/init"
	logs "github.com/DoRightt/pickfighter-server/pkg/logger"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:          "repo",
	Short:        "helps to communicate with application PostgreSQL database",
	Long:         ``,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(initSchemaCmd)
}

var initSchemaCmd = &cobra.Command{
	Use:   "init-db",
	Short: "DB initialization",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		db, err := psql.New(ctx)
		if err != nil {
			logs.Fatalf("Unable to connect PostgreSQL: %s", err)
		}
		defer db.GracefulShutdown()

		err = migrations.InitEventsSchema(ctx, db)
		if err != nil {
			logs.Fatalf("Error while initializing database : %v", err)
		}
		fmt.Println("Database successfully initalized")
	},
}
