package model

import (
	"time"
)

// RegisterRequest represents the data structure for handling user registration requests.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	TermsOk  bool   `json:"terms_ok"`

	Token string `json:"token"`
}

// AuthenticateRequest represents the data structure for user authentication requests.
type AuthenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	RememberMe bool   `json:"remember_me"`
	UserAgent  string `json:"user_agent"`
	IpAddress  string `json:"ip_address"`

	Subject   string   `json:"subject"`
	ExpiresIn int64    `json:"expires_in"`
	Audience  []string `json:"audience"`

	Method int `json:"method"`
}

// AuthenticateResult represents the result of a successful authentication.
type AuthenticateResult struct {
	UserId         int32     `json:"user_id" yaml:"user_id"`
	TokenId        string    `json:"token_id" yaml:"token_id"`
	Code           string    `json:"code" yaml:"code"`
	AccessToken    string    `json:"access_token" yaml:"access_token"`
	ExpirationTime time.Time `json:"expiration_time" yaml:"expiration_time"`
}

// UserCredentials represents user authentication credentials and related information.
type UserCredentials struct {
	UserId      int32     `json:"user_id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Salt        string    `json:"-"`
	Token       string    `json:"-"`
	TokenType   TokenType `json:"token_type"`
	TokenExpire int64     `json:"token_expire"`
	Active      bool      `json:"active"`
}

// UserCredentialsRequest represents a request for retrieving user authentication credentials.
type UserCredentialsRequest struct {
	UserId    int32     `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenType TokenType `json:"token_type"`
	IsActive  int32     `json:"is_active"`
}

// ResetPasswordRequest represents a request to initiate the password reset process.
type ResetPasswordRequest struct {
	Email string `json:"email"`
}

// RecoverPasswordRequest represents a request to recover the password after initiating the reset process.
type RecoverPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// ChangePasswordRequest represents a request to change the user's password.
type ChangePasswordRequest struct {
	OldPassword    string `json:"old_password"`
	NewPassword    string `json:"new_password"`
	RepeatPassword string `json:"repeat_password"`
}
