package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repoCmd)
}

var repoCmd = &cobra.Command{
	Use:          "repo",
	Short:        "helps to communicate with application PostgreSQL database",
	Long:         ``,
	SilenceUsage: true,
}
