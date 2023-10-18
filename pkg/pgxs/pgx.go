package pgxs

import (
	"fmt"

	"go.uber.org/zap"
)

type Config struct {
	DataDir  string
	DbUri    string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (c *Config) GetConnString() string {
	if len(c.DbUri) > 0 {
		return c.DbUri
	}

	connString := fmt.Sprintf("host=%s port=%s database=%s user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
	)

	return connString
}

type Repo struct {
	Logger *zap.Logger `json:"-" yaml:"-"`
	Config *Config     `json:"-" yaml:"-"`
}
