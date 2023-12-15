package internalErr

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

	Profile = 500

	Token        = 600
	TokenEmpty   = 601
	TokenExpired = 602

	JSON        = 700
	JSONDecoder = 701
	JSONEncoder = 702

	DB        = 800
	DBGetUser = 801

	Events            = 900
	EventsFightResult = 901

	Count         = 1000
	CountFighters = 1001
	CountEvents   = 1002

	Fighters = 1100
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
	case TxCommit:
		return &InternalError{Code: code, Message: "[Transaction] Failed to commit registration transaction"}
	case TxNotUnique:
		return &InternalError{Code: code, Message: "[Transaction] Value already exists"}
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
	case AuthFormPasswordsMismatch:
		return &InternalError{Code: code, Message: "[Auth]: Passwords mismatch"}
	case QueryParamsToken:
		return &InternalError{Code: code, Message: "[Query Params]: Query parameter 'token' should be specified"}
	case UserCredentials:
		return &InternalError{Code: code, Message: "[User Credentials]: Failed to get user credentials"}
	case UserCredentialsToken:
		return &InternalError{Code: code, Message: "[User Credentials]: User credentials with specified token does not exists"}
	case UserCredentialsIsNotActive:
		return &InternalError{Code: code, Message: "[User Credentials]: User is not activated"}
	case UserCredentialsReset:
		return &InternalError{Code: code, Message: "[User Credentials]: Failed to update user password"}
	case Profile:
		return &InternalError{Code: code, Message: "[Profile]: Failed to find user profile"}
	case Token:
		return &InternalError{Code: code, Message: "[Token]: Token unknown error"}
	case TokenEmpty:
		return &InternalError{Code: code, Message: "[Token]: Token is empty"}
	case TokenExpired:
		return &InternalError{Code: code, Message: "[Token]: Token expired, try to reset password"}
	case JSON:
		return &InternalError{Code: code, Message: "[JSON]: JSON unknown error"}
	case JSONDecoder:
		return &InternalError{Code: code, Message: "[JSON]: Decoder error"}
	case DBGetUser:
		return &InternalError{Code: code, Message: "[DB]: Failed to get user"}
	case Events:
		return &InternalError{Code: code, Message: "[Events]: Decode error"}
	case EventsFightResult:
		return &InternalError{Code: code, Message: "[Events]: Failed to set fight result"}
	case Count:
		return &InternalError{Code: code, Message: "[Count]: Failed to get items count"}
	case CountFighters:
		return &InternalError{Code: code, Message: "[Count]: Failed to get fighters count"}
	case CountEvents:
		return &InternalError{Code: code, Message: "[Count]: Failed to get events count"}
	case Fighters:
		return &InternalError{Code: code, Message: "[Fighters]: Failed to find fighters"}
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
