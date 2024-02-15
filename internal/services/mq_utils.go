package services

import (
	"context"
	"fmt"
	"log"
	"projects/fb-server/pkg/model"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// HandleEmailEvent processes different email events based on the provided EmailData.
// It uses the gomail package to send emails through an SMTP server. The email content
// and recipient details are determined by the event type, such as registration or
// password reset.
func (h *ApiHandler) HandleEmailEvent(ctx context.Context, data *model.EmailData) {
	d := gomail.NewDialer("smtp.gmail.com", 587, viper.GetString("mail.sender_address"), viper.GetString("mail.app_password"))
	host := viper.GetString("web.host")
	port := viper.GetString("web.port")

	var m *gomail.Message

	switch data.Subject {
	case model.EmailRegistration:
		m = getVerificationMessage(data, host, port)
	case model.EmailResetPassword:
		m = getPasswordRecoveryMessage(data, host, port)
	default:
		fmt.Println("Unexpected subject")
	}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Unable to send your email")
		log.Fatal(err)
	}
}

// getVerificationMessage generates a verification email message.
// It utilizes the provided EmailData and server information (host, port) to create
// a message with a verification link. The sender, recipient, subject, and body are
// set accordingly in the gomail.Message.
func getVerificationMessage(data *model.EmailData, host, port string) *gomail.Message {
	m := gomail.NewMessage()

	text := fmt.Sprintf("Hello, here is your verification link: %s:%s/register/confirm?token=%s", host, port, data.Token)

	m.SetHeader("From", viper.GetString("mail.sender_address"))
	m.SetHeader("To", data.Recipient.Email)
	m.SetHeader("Subject", "Please, Verify your email.")

	m.SetBody("text/plain", text)

	return m
}

// getPasswordRecoveryMessage generates a password recovery email message using the provided
// EmailData and server information (host, port). It constructs a message with a recovery link
// containing the host, port, and token. The email sender and recipient addresses, as well as
// the subject, are set in the message headers. The message body is a plain text representation
// containing the recovery link.
func getPasswordRecoveryMessage(data *model.EmailData, host, port string) *gomail.Message {
	m := gomail.NewMessage()

	text := fmt.Sprintf("Hello, here you can change your password: %s:%s/password/recover?token=%s", host, port, data.Token)

	m.SetHeader("From", viper.GetString("mail.sender_address"))
	m.SetHeader("To", data.Recipient.Email)
	m.SetHeader("Subject", "Please, Set a new password")

	m.SetBody("text/plain", text)

	return m
}
