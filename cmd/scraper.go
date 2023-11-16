package cmd

import (
	"projects/fb-server/pkg/scraper"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scrapeCmd)
}

var scrapeCmd = &cobra.Command{
	Use:              "scrape",
	Short:            "Run WEB Scraper",
	Long:             ``,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		scraper.Run()
	},
}
