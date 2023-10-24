package model

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	TermsOk  bool   `json:"terms_ok"`

	Token string `json:"token"`
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
