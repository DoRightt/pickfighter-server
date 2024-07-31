package internal

import (
	"fmt"
	"time"

	"fightbettr.com/fighters/pkg/version"
)

type DefaultMessagesList map[int]Error

const (
	Tx          int = 10
	TxCommit        = 11
	TxNotUnique     = 12
	TxUnknown       = 19

	Auth                      = 200
	AuthDecode                = 210
	AuthForm                  = 220
	AuthFormEmailEmpty        = 221
	AuthFormEmailInvalid      = 222
	AuthFormPasswordInvalid   = 223
	AuthFormPasswordWrong     = 224
	AuthFormPasswordsMismatch = 225

	QueryParams      = 300
	QueryParamsToken = 301

	UserCredentials            = 400
	UserCredentialsNotExists   = 401
	UserCredentialsToken       = 402
	UserCredentialsIsNotActive = 403
	UserCredentialsReset       = 404
	UserCredentialsCreate      = 405
	UserCredentialsUpdate      = 406

	Profile = 500

	Token        = 600
	TokenEmpty   = 601
	TokenExpired = 602

	JSON        = 700
	JSONDecoder = 701
	JSONEncoder = 702

	DB        = 800
	DBGetUser = 801

	Count         = 1000
	CountFighters = 1001
	Fighters      = 1002
)

var defaultErrors = DefaultMessagesList{
	Tx:                         Error{Code: Tx, Message: "[Transaction] Failed transaction"},
	TxCommit:                   Error{Code: TxCommit, Message: "[Transaction] Failed to commit registration transaction"},
	TxNotUnique:                Error{Code: TxNotUnique, Message: "[Transaction] Value already exists"},
	TxUnknown:                  Error{Code: TxUnknown, Message: "[Transaction] Failed transaction"},
	Auth:                       Error{Code: Auth, Message: "[Auth] Error"},
	AuthDecode:                 Error{Code: AuthDecode, Message: "[Auth]: Decode Error"},
	AuthForm:                   Error{Code: AuthForm, Message: "[Auth]: Form data is invalid"},
	AuthFormEmailEmpty:         Error{Code: AuthFormEmailEmpty, Message: "[Auth]: Email is empty"},
	AuthFormEmailInvalid:       Error{Code: AuthFormEmailInvalid, Message: "[Auth]: Email address is invalid"},
	AuthFormPasswordInvalid:    Error{Code: AuthFormPasswordInvalid, Message: "[Auth]: Password is empty or less than 6 symbols"},
	AuthFormPasswordWrong:      Error{Code: AuthFormPasswordWrong, Message: "[Auth]: Wrong Password"},
	AuthFormPasswordsMismatch:  Error{Code: AuthFormPasswordsMismatch, Message: "[Auth]: Passwords mismatch"},
	QueryParamsToken:           Error{Code: QueryParamsToken, Message: "[Query Params]: Query parameter 'token' should be specified"},
	UserCredentials:            Error{Code: UserCredentials, Message: "[User Credentials]: Failed to get user credentials"},
	UserCredentialsNotExists:   Error{Code: UserCredentialsNotExists, Message: "[User Credentials]: User with specified login credentials not exists"},
	UserCredentialsToken:       Error{Code: UserCredentialsToken, Message: "[User Credentials]: User credentials with specified token does not exist"},
	UserCredentialsIsNotActive: Error{Code: UserCredentialsIsNotActive, Message: "[User Credentials]: User is not activated"},
	UserCredentialsReset:       Error{Code: UserCredentialsReset, Message: "[User Credentials]: Failed to update user password"},
	UserCredentialsCreate:      Error{Code: UserCredentialsCreate, Message: "[User Credentials]: Failed to create user credentials"},
	UserCredentialsUpdate:      Error{Code: UserCredentialsUpdate, Message: "[User Credentials]: Failed to update user credentials"},
	Profile:                    Error{Code: Profile, Message: "[Profile]: Failed to find user profile"},
	Token:                      Error{Code: Token, Message: "[Token]: Token unknown error"},
	TokenEmpty:                 Error{Code: TokenEmpty, Message: "[Token]: Token is empty"},
	TokenExpired:               Error{Code: TokenExpired, Message: "[Token]: Token expired, try to reset password"},
	JSON:                       Error{Code: JSON, Message: "[JSON]: JSON unknown error"},
	JSONDecoder:                Error{Code: JSONDecoder, Message: "[JSON]: Decoder error"},
	DBGetUser:                  Error{Code: DBGetUser, Message: "[DB]: Failed to get user"},
	Count:                      Error{Code: Count, Message: "[Count]: Failed to get items count"},
	CountFighters:              Error{Code: CountFighters, Message: "[Count]: Failed to get fighters count"},
	Fighters:                   Error{Code: Fighters, Message: "[Fighters]: Failed to find fighters"},
}

var unknownError = Error{Code: 9999, Message: "Unknown Error"}

type Error struct {
	Code         int
	InternalCode int
	Message      string
	Timestamp    string
}

func (e *Error) GetCode() int {
	return int(e.Code)
}

func (e *Error) GetMessage() string {
	return e.Message
}

func NewDefault(code int, internal int) *Error {
	err, ok := defaultErrors[code]
	if !ok {
		err = unknownError
	}

	err.InternalCode = internal
	err.Timestamp = time.Now().Format(time.RFC1123)

	return &err
}

func New(code int, err error, internal int) *Error {
	return &Error{
		Code:         code,
		InternalCode: internal,
		Message:      err.Error(),
		Timestamp:    time.Now().Format(time.RFC1123),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"[ERROR]: %s. \n[ERROR CODE]: %d. \n[INTERNAL CODE]: %d. \nSERVICE: %s.\nTime: %s.\n",
		e.Message,
		e.Code,
		e.InternalCode,
		version.Name,
		e.Timestamp,
	)
}
