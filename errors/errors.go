package internalErr

const (
	Tx          int = 10
	TxCommit        = 11
	TxNotUnique     = 12
	TxUnknown       = 19

	Auth                    = 200
	AuthDecode              = 210
	AuthForm                = 220
	AuthFormEmailEmpty      = 221
	AuthFormEmailInvalid    = 222
	AuthFormPasswordInvalid = 223
	AuthFormPasswordWrong   = 224

	QueryParams      = 300
	QueryParamsToken = 301

	UserCredentials            = 400
	UserCredentialsNotExists   = 401
	UserCredentialsToken       = 402
	UserCredentialsIsNotActive = 403

	Token        = 500
	TokenExpired = 501

	JSON        = 600
	JSONDecoder = 601
	JSONEncoder = 602

	DB        = 700
	DBGetUser = 701
)

type InternalError struct {
	Code    int
	Message string
}

func (e *InternalError) GetCode() int {
	return int(e.Code)
}

func New(code int) *InternalError {
	switch code {
	case Tx:
		return &InternalError{Code: code, Message: "[Transaction] Failed transaction"}
	case TxNotUnique:
		return &InternalError{Code: code, Message: "[Transaction] Email already exists"}
	case TxUnknown:
		return &InternalError{Code: code, Message: "[Transaction] Failed transaction"}
	case Auth:
		return &InternalError{Code: code, Message: "[Auth] Error"}
	case AuthDecode:
		return &InternalError{Code: code, Message: "[Auth]: Decode Error"}
	case AuthForm:
		return &InternalError{Code: code, Message: "[Auth]: Form data is invalid"}
	case AuthFormEmailEmpty:
		return &InternalError{Code: code, Message: "[Auth]: Email is empty"}
	case AuthFormEmailInvalid:
		return &InternalError{Code: code, Message: "[Auth]: Email address is invalid"}
	case AuthFormPasswordInvalid:
		return &InternalError{Code: code, Message: "[Auth]: Password is empty or less than 6 symbols"}
	case AuthFormPasswordWrong:
		return &InternalError{Code: code, Message: "[Auth]: Wrong Password"}
	case QueryParamsToken:
		return &InternalError{Code: code, Message: "[Query Params]: Query parameter 'token' should be specified"}
	case UserCredentials:
		return &InternalError{Code: code, Message: "[User Credentials]: Failed to get user credentials"}
	case UserCredentialsToken:
		return &InternalError{Code: code, Message: "[User Credentials]: User credentials with specified token does not exists"}
	case UserCredentialsIsNotActive:
		return &InternalError{Code: code, Message: "[User Credentials]: User is not activated"}
	case TokenExpired:
		return &InternalError{Code: code, Message: "[Token]: Token expired, try to reset password"}
	case JSON:
		return &InternalError{Code: code, Message: "[JSON]: JSON unknown error"}
	case JSONDecoder:
		return &InternalError{Code: code, Message: "[JSON]: Decoder error"}
	case DBGetUser:
		return &InternalError{Code: code, Message: "[DB]: Failed to get user"}
	default:
		return &InternalError{
			Code:    1001,
			Message: "Unknown error",
		}

	}
}

func (e *InternalError) Error() string {
	return e.Message
}
