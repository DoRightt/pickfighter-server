package cmd

import (
	"fmt"
	"log"
	"os"
	lg "projects/fb-server/logger"
	"projects/fb-server/pkg/version"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgPath string
	logger  *zap.SugaredLogger
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
		logger.Fatal(err.Error())
	}
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("Error loading .env file")
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Config file path (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("name", version.Name, "Application name label")

	bindViperPersistentFlag(rootCmd, "config_path", "config")
	bindViperPersistentFlag(rootCmd, "app.name", "name")

	rootCmd.Flags().BoolP("version", "v", false, "Shows app version")

	initZapLogger()
}

func initZapLogger() {
	logLevel := "info"
	logFilePath := "logger/logs/log.json"

	if err := lg.Initialize(logLevel, logFilePath); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	logger = lg.Get().Sugar()
}

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
	}
}

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

	// postgres
	viper.SetDefault("postgres.url", os.Getenv("POSTGRES_URL"))
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.name", "postgres")
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "")
}

func bindViperFlag(cmd *cobra.Command, viperVal, flagName string) {
	if err := viper.BindPFlag(viperVal, cmd.Flags().Lookup(flagName)); err != nil {
		log.Printf("Failed to bind viper flag: %s", err)
	}
}

func bindViperPersistentFlag(cmd *cobra.Command, viperVal, flagName string) {
	if err := viper.BindPFlag(viperVal, cmd.PersistentFlags().Lookup(flagName)); err != nil {
		log.Printf("Failed to bind viper flag: %s", err)
	}
}
