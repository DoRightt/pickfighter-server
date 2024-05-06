package cfg

import (
	"fightbettr.com/fb-server/pkg/pgxs"

	"github.com/spf13/viper"
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
