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
	case AuthFormEmailEmpty:
		return &InternalError{Code: code, Message: "[Auth]: Email is empty"}
	case AuthFormEmailInvalid:
		return &InternalError{Code: code, Message: "[Auth]: Email address is invalid"}
	case AuthFormPasswordInvalid:
		return &InternalError{Code: code, Message: "[Auth]: Password is empty or less than 6 symbols"}
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
