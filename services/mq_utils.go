package services

import (
	"context"
	"fmt"
	"log"
	"projects/fb-server/pkg/model"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func (h *ApiHandler) HandleEmailEvent(ctx context.Context, data *model.EmailData) {
	d := gomail.NewDialer("smtp.gmail.com", 587, viper.GetString("mail.sender_address"), viper.GetString("mail.app_password"))
	m := gomail.NewMessage()
	host := viper.GetString("web.host")
	port := viper.GetString("web.port")
	text := fmt.Sprintf("Hello, here is your verification link: %s:%s/register/confirm?token=%s", host, port, data.Token)

	m.SetHeader("From", viper.GetString("mail.sender_address"))
	m.SetHeader("To", data.Recipient.Email)
	m.SetHeader("Subject", "Please, Verify your email.")

	m.SetBody("text/plain", text)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Unable to send your email")
		log.Fatal(err)
	}
}
