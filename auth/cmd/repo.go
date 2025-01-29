package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"pickfighter.com/auth/internal/repository/psql"
	logs "pickfighter.com/pkg/logger"
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
