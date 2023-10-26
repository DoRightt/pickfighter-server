package services

import (
	"context"
	"fmt"
	"log"
	"projects/fb-server/pkg/model"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

func (h *ApiHandler) HandleEmailEvent(ctx context.Context, data *model.EmailData) {
	apiKey := viper.GetString("sendgrid.api_key")
	senderName := viper.GetString("sendgrid.sender_name")
	senderAddr := viper.GetString("sendgrid.sender_address")

	from := mail.NewEmail(senderName, senderAddr)
	subject := "Please, Verificate your email."
	to := mail.NewEmail(data.Recipient.Name, data.Recipient.Email)
	plainTextContent := "Here is your AMAZING email!"
	htmlContent := "Here is your <strong>AMAZING</strong> email!"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		fmt.Println("Unable to send your email")
		log.Fatal(err)
	}

	statusCode := response.StatusCode
	if statusCode == 200 || statusCode == 201 || statusCode == 202 {
		fmt.Println("Email sent!")
	}
}
