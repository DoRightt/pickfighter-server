package cmd

import (
	"fmt"
	"log"
	"time"

	"fightbettr.com/scraper/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgPath string
)

// rootCmd is the main Cobra command representing the root of the scraper service.
var rootCmd = &cobra.Command{
	Use:   "UFC-Scraper",
	Short: "UFC fighters scraper",
	RunE: func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Println("Dev version", version.DevVersion)
			return nil
		}

		return cmd.Usage()
	},
}

// Execute runs the root command for the scraper service.
// It executes the necessary logic for the command-line interface,
// handling errors and logging them if they occur.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Config file path (default is ./configs/proxy.yaml)")
	rootCmd.PersistentFlags().Bool("proxy", false, "Run with proxy")
	rootCmd.PersistentFlags().Bool("add", false, "Add results to previus fighters collection")
	rootCmd.PersistentFlags().Int("start", 0, "start page")

	bindViperPersistentFlag(rootCmd, "config_path", "config")
	bindViperPersistentFlag(rootCmd, "proxy", "proxy")
	bindViperPersistentFlag(rootCmd, "add", "add")
	bindViperPersistentFlag(rootCmd, "start", "start")
}

// initConfig initializes the service configuration.
// It sets default values, reads from environment variables, and reads from a config file if present.
func initConfig() {
	setConfigDefaults()

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
		viper.SetConfigName("proxy")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setConfigDefaults() {
	// app defaults
	viper.SetDefault("app.env", "dev")
	viper.SetDefault("app.name", version.Name)
	viper.SetDefault("app.version", version.DevVersion)
	viper.SetDefault("app.run_date", time.Unix(version.RunDate, 0).Format(time.RFC1123))
}

// bindViperPersistentFlag binds a Viper configuration flag to a persistent Cobra command flag.
func bindViperPersistentFlag(cmd *cobra.Command, viperVal, flagName string) {
	if err := viper.BindPFlag(viperVal, cmd.PersistentFlags().Lookup(flagName)); err != nil {
		log.Printf("Failed to bind viper flag: %s", err)
	}
}
