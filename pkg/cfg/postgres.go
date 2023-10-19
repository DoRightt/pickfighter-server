package cfg

import (
	"projects/fb-server/pkg/pgxs"

	"github.com/spf13/viper"
)

func ViperPostgres() *pgxs.Config {
	return &pgxs.Config{
		DataDir:  viper.GetString("postgres.data_dir"),
		DbUri:    viper.GetString("postgres.url"),
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
		Name:     viper.GetString("postgres.name"),
		User:     viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
	}
}
