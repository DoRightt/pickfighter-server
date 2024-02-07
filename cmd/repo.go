package cmd

import (
	"context"
	"fmt"
	lg "projects/fb-server/pkg/logger"
	"projects/fb-server/pkg/cfg"
	"projects/fb-server/pkg/pgxs"
	"projects/fb-server/internal/scraper"
	"sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repoCmd)
	rootCmd.AddCommand(updateRosterCmd)
	rootCmd.AddCommand(clearTableCmd)

	updateRosterCmd.Flags().Bool("scrape", false, "Scrape new data before update")
}

var repoCmd = &cobra.Command{
	Use:          "repo",
	Short:        "helps to communicate with application PostgreSQL database",
	Long:         ``,
	SilenceUsage: true,
}

var updateRosterCmd = &cobra.Command{
	Use:              "update",
	Short:            "Updates fighters table by json list",
	Long:             ``,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		scrapeFlag, _ := cmd.Flags().GetBool("scrape")

		if scrapeFlag {
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				scraper.Run()
			}()

			wg.Wait()
		}

		l := lg.NewSugared()

		fighters, err := scraper.ReadFighterData()
		if err != nil {
			logger.Fatalf("Error while reading figheter data: %s", err)
		}

		scraper.WriteFighterData(ctx, l, fighters)
	},
}

var clearTableCmd = &cobra.Command{
	Use:   "clear-table [table_name]",
	Short: "Delete all records from the specified table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		t := args[0]

		db, err := pgxs.NewPool(ctx, logger, cfg.ViperPostgres())
		if err != nil {
			logger.Fatalf("Unable to connect postgresql: %s", err)
		}

		err = db.DeleteRecords(ctx, t)
		if err != nil {
			logger.Fatalf("Error deleting records: %s", err)
		}
		fmt.Printf("All records from table '%s' deleted successfully\n", t)
	},
}
