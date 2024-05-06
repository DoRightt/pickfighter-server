package model

// Email types used in the application.
const (
	EmailRegistration  = "registration"
	EmailResetPassword = "reset_password"
)

// EmailAddrSpec represents an email address with an optional name.
type EmailAddrSpec struct {
	Email string `json:"email" yaml:"email"`
	Name  string `json:"name" yaml:"name"`
}

// EmailData represents the data structure for sending emails in the application.
type EmailData struct {
	Sender    EmailAddrSpec `json:"sender" yaml:"sender"`
	Recipient EmailAddrSpec `json:"recipient" yaml:"recipient"`
	Subject   string        `json:"subject" yaml:"subject"`
	Token     string        `json:"token" yaml:"token"`
	Url       string        `json:"url" yaml:"url"`
}
