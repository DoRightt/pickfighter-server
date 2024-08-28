package internal

import (
	"fmt"
	"time"

	"pickfighter.com/auth/pkg/version"
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
)

var defaultErrors = DefaultMessagesList{
	Tx:                         Error{ErrCode: Tx, Message: "[Transaction] Failed transaction"},
	TxCommit:                   Error{ErrCode: TxCommit, Message: "[Transaction] Failed to commit registration transaction"},
	TxNotUnique:                Error{ErrCode: TxNotUnique, Message: "[Transaction] Value already exists"},
	TxUnknown:                  Error{ErrCode: TxUnknown, Message: "[Transaction] Failed transaction"},
	Auth:                       Error{ErrCode: Auth, Message: "[Auth] Error"},
	AuthDecode:                 Error{ErrCode: AuthDecode, Message: "[Auth]: Decode Error"},
	AuthForm:                   Error{ErrCode: AuthForm, Message: "[Auth]: Form data is invalid"},
	AuthFormEmailEmpty:         Error{ErrCode: AuthFormEmailEmpty, Message: "[Auth]: Email is empty"},
	AuthFormEmailInvalid:       Error{ErrCode: AuthFormEmailInvalid, Message: "[Auth]: Email address is invalid"},
	AuthFormPasswordInvalid:    Error{ErrCode: AuthFormPasswordInvalid, Message: "[Auth]: Password is empty or less than 6 symbols"},
	AuthFormPasswordWrong:      Error{ErrCode: AuthFormPasswordWrong, Message: "[Auth]: Wrong Password"},
	AuthFormPasswordsMismatch:  Error{ErrCode: AuthFormPasswordsMismatch, Message: "[Auth]: Passwords mismatch"},
	QueryParamsToken:           Error{ErrCode: QueryParamsToken, Message: "[Query Params]: Query parameter 'token' should be specified"},
	UserCredentials:            Error{ErrCode: UserCredentials, Message: "[User Credentials]: Failed to get user credentials"},
	UserCredentialsNotExists:   Error{ErrCode: UserCredentialsNotExists, Message: "[User Credentials]: User with specified login credentials not exists"},
	UserCredentialsToken:       Error{ErrCode: UserCredentialsToken, Message: "[User Credentials]: User credentials with specified token does not exist"},
	UserCredentialsIsNotActive: Error{ErrCode: UserCredentialsIsNotActive, Message: "[User Credentials]: User is not activated"},
	UserCredentialsReset:       Error{ErrCode: UserCredentialsReset, Message: "[User Credentials]: Failed to update user password"},
	UserCredentialsCreate:      Error{ErrCode: UserCredentialsCreate, Message: "[User Credentials]: Failed to create user credentials"},
	UserCredentialsUpdate:      Error{ErrCode: UserCredentialsUpdate, Message: "[User Credentials]: Failed to update user credentials"},
	Profile:                    Error{ErrCode: Profile, Message: "[Profile]: Failed to find user profile"},
	Token:                      Error{ErrCode: Token, Message: "[Token]: Token unknown error"},
	TokenEmpty:                 Error{ErrCode: TokenEmpty, Message: "[Token]: Token is empty"},
	TokenExpired:               Error{ErrCode: TokenExpired, Message: "[Token]: Token expired, try to reset password"},
	JSON:                       Error{ErrCode: JSON, Message: "[JSON]: JSON unknown error"},
	JSONDecoder:                Error{ErrCode: JSONDecoder, Message: "[JSON]: Decoder error"},
	DBGetUser:                  Error{ErrCode: DBGetUser, Message: "[DB]: Failed to get user"},
}

var unknownError = Error{ErrCode: 9999, Message: "Unknown Error"}

type Error struct {
	ErrCode      int
	InternalCode int
	Message      string
	Timestamp    any
}

func (e *Error) GetCode() int {
	return int(e.ErrCode)
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
		ErrCode:      code,
		InternalCode: internal,
		Message:      err.Error(),
		Timestamp:    time.Now().Format(time.RFC1123),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"[ERROR]: %s. \n[ERROR CODE]: %d. \n[INTERNAL CODE]: %d. \nSERVICE: %s.\nTime: %t.\n",
		e.Message,
		e.ErrCode,
		e.InternalCode,
		version.Name,
		e.Timestamp,
	)
}
