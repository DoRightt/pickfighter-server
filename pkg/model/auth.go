package model

import "time"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	TermsOk  bool   `json:"terms_ok"`

	Token string `json:"token"`
}

type LoginRequest struct {
}

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

type AuthenticateResult struct {
	UserId         int32     `json:"user_id" yaml:"user_id"`
	TokenId        string    `json:"token_id" yaml:"token_id"`
	Code           string    `json:"code" yaml:"code"`
	AccessToken    string    `json:"access_token" yaml:"access_token"`
	ExpirationTime time.Time `json:"expiration_time" yaml:"expiration_time"`
}

type UserCredentials struct {
	UserId      int32  `json:"user_id"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	Salt        string `json:"-"`
	Token       string `json:"-"`
	TokenType   string `json:"token_type"`
	TokenExpire int64  `json:"token_expire"`
	Active      bool   `json:"active"`
}

type UserCredentialsRequest struct {
	UserId    int32  `json:"user_id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	IsActive  int32  `json:"is_active"`
}
