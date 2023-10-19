package services

import (
	"context"
	"projects/fb-server/pkg/httplib"

	"github.com/spf13/viper"
)

func (h *ApiHandler) RunHTTPServer(ctx context.Context) error {
	httplib.SetCookieName(viper.GetString("auth.cookie_name"))
	
}
