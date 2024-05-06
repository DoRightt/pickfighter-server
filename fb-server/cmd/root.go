// Package cmd provides the command-line interface for the FB-Server application.
// It includes commands for running the server, displaying version information, and more.
package cmd

import (
	"fmt"
	"log"
	"os"
	lg "fightbettr.com/fb-server/pkg/logger"
	"fightbettr.com/fb-server/pkg/version"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

var (
	cfgPath string
	logger  lg.FbLogger
)

// rootCmd is the main Cobra command representing the root of the CLI application.
var rootCmd = &cobra.Command{
	Use:   "excelsior",
	Short: "FB-Server, CLI App",
	RunE: func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Println("Dev version", version.DevVersion)
			fmt.Println("Git version", version.GitVersion)
			return nil
		}
		return cmd.Usage()
	},
}

// Execute runs the root command for the fb-server application.
// It executes the necessary logic for the command-line interface,
// handling errors and logging them if they occur.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("Error loading .env file")
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Config file path (default is ./config.yaml)")
	rootCmd.PersistentFlags().Bool("log_json", false, "Enable JSON formatted logs output")
	rootCmd.PersistentFlags().Int("log_level", int(zapcore.DebugLevel), "Log level")
	rootCmd.PersistentFlags().String("name", version.Name, "Application name label")

	bindViperPersistentFlag(rootCmd, "config_path", "config")
	bindViperPersistentFlag(rootCmd, "app.name", "name")
	bindViperPersistentFlag(rootCmd, "log_json", "log_json")
	bindViperPersistentFlag(rootCmd, "log_level", "log_level")

	rootCmd.Flags().BoolP("version", "v", false, "Shows app version")
}

// initZapLogger initializes the zap logger.
func initZapLogger() {
	logger = lg.NewSugared()
}

// initConfig initializes the application configuration.
// It sets default values, reads from environment variables, and reads from a config file if present.
func initConfig() {
	setConfigDefaults()

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		initZapLogger()
	}
}

// setConfigDefaults sets default values for various configuration options.
func setConfigDefaults() {
	// app defaults
	viper.SetDefault("app.env", "dev")
	viper.SetDefault("app.name", version.Name)
	viper.SetDefault("app.version", version.DevVersion)
	viper.SetDefault("app.build_date", version.BuildDate)
	viper.SetDefault("app.run_date", time.Unix(version.RunDate, 0).Format(time.RFC1123))

	// http server
	viper.SetDefault("http.addr", "127.0.0.1:9090")
	viper.SetDefault("http.ssl.enabled", false)

	// auth config
	viper.SetDefault("auth.cookie_name", "fb_api_token")
	viper.SetDefault("auth.jwt.cert", "")
	viper.SetDefault("auth.jwt.key", "")

	// postgres
	viper.SetDefault("postgres.main.url", os.Getenv("POSTGRES_URL"))
	viper.SetDefault("postgres.main.host", "localhost")
	viper.SetDefault("postgres.main.port", "5432")
	viper.SetDefault("postgres.main.name", "postgres")
	viper.SetDefault("postgres.main.user", "postgres")

	// email
	viper.SetDefault("mail.sender_address", os.Getenv("SEND_FROM_ADDRESS"))
	viper.SetDefault("mail.sender_name", os.Getenv("SEND_FROM_NAME"))
	viper.SetDefault("mail.app_password", os.Getenv("MAIL_APP_PASSWORD"))

	// web
	viper.SetDefault("web.host", os.Getenv("FRONT_HOST"))
	viper.SetDefault("web.port", os.Getenv("FRONT_PORT"))
}

// bindViperFlag binds a Viper configuration flag to a Cobra command flag.
func bindViperFlag(cmd *cobra.Command, viperVal, flagName string) {
	if err := viper.BindPFlag(viperVal, cmd.Flags().Lookup(flagName)); err != nil {
		log.Printf("Failed to bind viper flag: %s", err)
	}
}

// bindViperPersistentFlag binds a Viper configuration flag to a persistent Cobra command flag.
func bindViperPersistentFlag(cmd *cobra.Command, viperVal, flagName string) {
	if err := viper.BindPFlag(viperVal, cmd.PersistentFlags().Lookup(flagName)); err != nil {
		log.Printf("Failed to bind viper flag: %s", err)
	}
}
