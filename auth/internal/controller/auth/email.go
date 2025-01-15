package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/spf13/viper"
	"pickfighter.com/auth/pkg/model"
)

// HandleEmailEvent processes different email events based on the provided EmailData.
// It uses the mailgun package to send emails. The email content
// and recipient details are determined by the event type, such as registration or
// password reset.
func (c *Controller) HandleEmailEvent(data *model.EmailData) {
	mg := mailgun.NewMailgun(viper.GetString("mail.mailgun_domain"), viper.GetString("mail.mailgun_api"))

	host := viper.GetString("web.host")
	port := viper.GetString("web.port")

	var message *mailgun.Message

	switch data.Subject {
	case model.EmailRegistration:
		message = getVerificationMessage(data, host, port)
	case model.EmailResetPassword:
		message = getPasswordRecoveryMessage(data, host, port)
	default:
		fmt.Println("Unexpected subject")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		fmt.Println("Unable to send your email")
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

// getVerificationMessage generates a verification email message.
// It utilizes the provided EmailData and server information (host, port) to create
// a message with a verification link. The sender, recipient, subject, and body are
// set accordingly in the mailgun.Message.
func getVerificationMessage(data *model.EmailData, host, port string) *mailgun.Message {
	sender := viper.GetString("mail.sender_address")
	subject := "Please, Verify your email."
	recipient := data.Recipient.Email
	body := fmt.Sprintf("Hello, here is your verification link: %s:%s/register/confirm?token=%s", host, port, data.Token)

	message := mailgun.NewMessage(sender, subject, body, recipient)

	return message
}

// getPasswordRecoveryMessage generates a password recovery email message using the provided
// EmailData and server information (host, port). It constructs a message with a recovery link
// containing the host, port, and token. The email sender and recipient addresses, as well as
// the subject, are set in the message headers. The message body is a plain text representation
// containing the recovery link.
func getPasswordRecoveryMessage(data *model.EmailData, host, port string) *mailgun.Message {
	sender := viper.GetString("mail.sender_address")
	subject := "Please, Set a new password"
	recipient := data.Recipient.Email
	body := fmt.Sprintf("Hello, here you can change your password: %s:%s/password/recover?token=%s", host, port, data.Token)

	message := mailgun.NewMessage(sender, subject, body, recipient)

	return message
}
