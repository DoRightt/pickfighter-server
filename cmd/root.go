package cmd

import (
	"errors"
	"fmt"
	"os"
	"projects/fb-server/logger"
	"projects/fb-server/pkg/version"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgPath string
	lg      *zap.Logger
)

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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		lg.Fatal(err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Config file path (default is $HOME/go/src/projects/fb-server/config.yaml)")

	rootCmd.Flags().BoolP("version", "v", false, "Shows app version")

	initZapLogger()
}

func initZapLogger() {
	logLevel := "info"
	logFilePath := "logger/logs/log.json"

	if err := logger.Initialize(logLevel, logFilePath); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	lg = logger.Get()

	lg.Info("This is an info log message")
	lg.Error("This is an error log message", zap.Error(errors.New("test")))
}

func initConfig() {
	setConfigDefaults()
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		viper.AddConfigPath(os.Getenv("HOME"))
		viper.AddConfigPath(os.Getenv("ENGINE_APP_HOME"))
		viper.SetConfigName("config")
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
	viper.SetDefault("app.build_date", version.BuildDate)
	viper.SetDefault("app.run_date", time.Unix(version.RunDate, 0).Format(time.RFC1123))

	// auth config
	viper.SetDefault("auth.cookie_name", "fb_api_token")

	// postgres
	viper.SetDefault("postgres.url", os.Getenv("POSTGRES_URL"))
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.name", "postgres")
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "")
}
