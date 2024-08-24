package cmd

import (
	"context"

	"pickfighter.com/fighters/pkg/cfg"
	"github.com/spf13/cobra"
)

func init() {
	repoCmd.AddCommand(clearTableCmd)
}

// clearTableCmd represents the clear-table command. It is used to delete all records from the fb_fighters and fb_fighter_stats tables.
// Expects one argument - the table name.
var clearTableCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all records from the fb_fighters and fb_fighter_stats tables",
	Args:  cobra.ExactArgs(1),
	Run:   runClearTable,
}

// runClearTable is the function executed when the clear-table command is run.
// It deletes all records from the fb_fighters and fb_fighter_stats tables.
func runClearTable(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := cfg.ViperPostgres()
	DeleteFighterData(ctx, cfg)
}
