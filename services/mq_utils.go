package services

import (
	"context"
	"projects/fb-server/pkg/model"

	"github.com/spf13/viper"
)

func (h *ApiHandler) HandleEmailEvent(ctx context.Context, data *model.EmailData) {
	apiKey := viper.GetString("sendgrid.api_key")

}
