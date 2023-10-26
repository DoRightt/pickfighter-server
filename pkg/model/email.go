package model

const (
	EmailRegistration  = "registration"
	EmailResetPassword = "reset_password"
)

type EmailAddrSpec struct {
	Email string `json:"email" yaml:"email"`
	Name  string `json:"name" yaml:"name"`
}

type EmailData struct {
	Sender    EmailAddrSpec `yaml:"sender" yaml:"sender"`
	Recipient EmailAddrSpec `json:"recipient" yaml:"recipient"`
	Subject   string        `json:"subject" yaml:"subject"`
	Token     string        `json:"token" yaml:"token"`
	Url       string        `json:"url" yaml:"url"`
}
