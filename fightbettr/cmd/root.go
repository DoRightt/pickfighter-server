package cmd

import (
	"fmt"

	lg "fightbettr.com/fightbettr/pkg/logger"
	"fightbettr.com/fighters/pkg/version"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	logger  lg.FbLogger
)

// rootCmd is the main Cobra command representing the root of the Fightbettr service.
var rootCmd = &cobra.Command{
	Use:   "Fightbettr Service",
	Short: "This CLI works with data to manage and redirect it",
	RunE: func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")

		if showVersion {
			fmt.Println("Dev version", version.DevVersion)
			return nil
		}

		return cmd.Usage()
	},
}

// Execute runs the root command for the Fightbettr service.
// It executes the necessary logic for the command-line interface,
// handling errors and logging them if they occur.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
