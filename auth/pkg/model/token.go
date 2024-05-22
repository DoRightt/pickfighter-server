package model

// TokenType represents the type for distinguishing token types in the application.
type TokenType string

// Constants for various token types.
const (
	TokenConfirmation  TokenType = "user_registration_confirmation"
	TokenResetPassword TokenType = "reset_password"
	TokenSetPassword   TokenType = "set_password"
)
