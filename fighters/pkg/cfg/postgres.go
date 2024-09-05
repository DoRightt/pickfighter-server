package cfg

import (
	"fmt"

	"github.com/spf13/viper"
	"pickfighter.com/pkg/pgxs"
)

// ViperPostgres returns the structure that returns pgxs.Cofnig based on values from viper
func ViperPostgres() *pgxs.Config {
	return &pgxs.Config{
		DataDir:  viper.GetString("postgres.main.data_dir"),
		DbUri:    viper.GetString("postgres.main.url"),
		Host:     viper.GetString("postgres.main.host"),
		Port:     viper.GetString("postgres.main.port"),
		Name:     viper.GetString("postgres.main.name"),
		User:     viper.GetString("postgres.main.user"),
		Password: viper.GetString("postgres.main.password"),
	}
}

func ViperTestPostgres() *pgxs.Config {
	config := &pgxs.Config{
		DataDir:  viper.GetString("postgres.test.data_dir"),
		DbUri:    viper.GetString("postgres.test.url"),
		Host:     viper.GetString("postgres.test.host"),
		Port:     viper.GetString("postgres.test.port"),
		Name:     viper.GetString("postgres.test.name"),
		User:     viper.GetString("postgres.test.user"),
		Password: viper.GetString("postgres.test.password"),
	}

	fmt.Printf("Config loaded:\nDataDir: %s\nDbUri: %s\nHost: %s\nPort: %s\nName: %s\nUser: %s\nPassword: %s\n",
		config.DataDir, config.DbUri, config.Host, config.Port, config.Name, config.User, config.Password)
	return config
}
