package cmd

import "github.com/spf13/cobra"

var repoCmd = &cobra.Command{
	Use:          "repo",
	Short:        "helps to communicate with application PostgreSQL database",
	Long:         ``,
	SilenceUsage: true,
}

func init() {
	
}