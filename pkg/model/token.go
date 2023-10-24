package model

type TokenType string

const (
	TokenConfirmation  = "user_registration_confirmation"
	TokenResetPassword = "reset_password"
	TokenSetPassword   = "set_password"
)