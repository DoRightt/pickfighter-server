package cmd

// func init() {
// rootCmd.AddCommand(repoCmd)
// rootCmd.AddCommand(updateRosterCmd)
// rootCmd.AddCommand(clearTableCmd)

// updateRosterCmd.Flags().Bool("scrape", false, "Scrape new data before update")
// }

// var repoCmd = &cobra.Command{
// 	Use:          "repo",
// 	Short:        "helps to communicate with application PostgreSQL database",
// 	Long:         ``,
// 	SilenceUsage: true,
// }

// updateRosterCmd represents the update command. It is used to update the fighters table using a JSON list.
// var updateRosterCmd = &cobra.Command{
// 	Use:              "update",
// 	Short:            "Updates fighters table by json list",
// 	Long:             ``,
// 	TraverseChildren: true,
// 	Run:              runUpdate,
// }

// clearTableCmd represents the clear-table command. It is used to delete all records from the specified table.
// Expects one argument - the table name.
// var clearTableCmd = &cobra.Command{
// 	Use:   "clear-table [table_name]",
// 	Short: "Delete all records from the specified table",
// 	Args:  cobra.ExactArgs(1),
// 	Run:   runClearTable,
// }

// runUpdate is the function executed when the update command is run.
// If the "scrape" argument is passed, the command will first run the scrapper to update the data and then update the table.
// If this argument is not passed, the table will be updated from the existing JSON file.
// func runUpdate(cmd *cobra.Command, args []string) {
// 	ctx := context.Background()
// 	scrapeFlag, _ := cmd.Flags().GetBool("scrape")

// 	if scrapeFlag {
// 		var wg sync.WaitGroup
// 		wg.Add(1)

// 		go func() {
// 			defer wg.Done()
// 			scraper.Run()
// 		}()

// 		wg.Wait()
// 	}

// 	l := lg.NewSugared()

// 	fighters, err := scraper.ReadFighterData()
// 	if err != nil {
// 		logger.Fatalf("Error while reading figheter data: %s", err)
// 	}

// 	scraper.WriteFighterData(ctx, l, fighters)
// }

// runClearTable is the function executed when the clear-table command is run.
// It deletes all records from the specified table in database.
// func runClearTable(cmd *cobra.Command, args []string) {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	t := args[0]

// 	db, err := pgxs.NewPool(ctx, logger, cfg.ViperPostgres())
// 	if err != nil {
// 		logger.Fatalf("Unable to connect postgresql: %s", err)
// 	}

// 	err = db.DeleteRecords(ctx, t)
// 	if err != nil {
// 		logger.Fatalf("Error deleting records: %s", err)
// 	}
// 	fmt.Printf("All records from table '%s' deleted successfully\n", t)
// }
