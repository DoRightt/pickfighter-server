package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/DoRightt/pickfighter-server/auth/internal/repository/psql"
	migrations "github.com/DoRightt/pickfighter-server/auth/migrations/init"
	logs "github.com/DoRightt/pickfighter-server/pkg/logger"
)

var repoCmd = &cobra.Command{
	Use:          "repo",
	Short:        "helps to communicate with application PostgreSQL database",
	Long:         ``,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoInitialRootUserCmd)
	repoCmd.AddCommand(initSchemaCmd)
}

var repoInitialRootUserCmd = &cobra.Command{
	Use:              "init-root-user",
	Short:            "Creates root user data",
	Long:             ``,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		db, err := psql.New(ctx)
		if err != nil {
			logs.Fatalf("Unable to connect postgresql: %s", err)
		}
		defer db.GracefulShutdown()

		return db.InitRootUser(ctx)
	},
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

		err = migrations.InitAuthSchema(ctx, db)
		if err != nil {
			logs.Fatalf("Error while initializing database : %v", err)
		}
		fmt.Println("Database successfully initalized")

		db.InitRootUser(ctx)
	},
}
