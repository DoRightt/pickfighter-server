package services

import (
	"projects/fb-server/pkg/model"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/gomail.v2"
)

func TestHandleEmailEvent(t *testing.T) {
	// TODO
}

func TestEmailMessage(t *testing.T) {
	emailData := &model.EmailData{
		Token: "testToken",
		Recipient: model.EmailAddrSpec{
			Email: "test@example.com",
		},
	}
	host := "example.com"
	port := "8080"

	tests := []struct {
		name     string
		message  *gomail.Message
		expected string
	}{
		{
			name:     "verification message",
			message:  getVerificationMessage(emailData, host, port),
			expected: "Please, Verify your email.",
		},
		{
			name:     "verification message",
			message:  getPasswordRecoveryMessage(emailData, host, port),
			expected: "Please, Set a new password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NotNil(t, tc.message)

			assert.Equal(t, viper.GetString("mail.sender_address"), tc.message.GetHeader("From")[0], "Header must match")
			assert.Equal(t, emailData.Recipient.Email, tc.message.GetHeader("To")[0], "Header must match")
			assert.Equal(t, tc.expected, tc.message.GetHeader("Subject")[0], "Header must match")

			tc.message = nil
		})
	}
}
