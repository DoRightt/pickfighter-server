package cmd

import (
	"context"

	"fightbettr.com/fighters/pkg/cfg"
	logs "fightbettr.com/pkg/logger"
	"github.com/spf13/cobra"
)

func init() {
	repoCmd.AddCommand(updateRosterCmd)
}

// updateRosterCmd represents the update command. It is used to update the fighters table using a JSON list.
var updateRosterCmd = &cobra.Command{
	Use:              "update",
	Short:            "Updates fighters table by json list",
	Long:             ``,
	TraverseChildren: true,
	Run:              runUpdate,
}

// runUpdate is the function executed when the update command is run.
// The table will be updated from the existing JSON file.
func runUpdate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	fighters, err := ReadFighterData()
	if err != nil {
		logs.Fatalf("Error while reading figheter data: %s", err)
	}

	cfg := cfg.ViperPostgres()
	WriteFighterData(ctx, fighters, cfg)
}
