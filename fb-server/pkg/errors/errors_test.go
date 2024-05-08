package internalErr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCode(t *testing.T) {
	tests := []struct {
		err    *InternalError
		expect int
	}{
		{&InternalError{Code: Tx, Message: "Transaction error"}, Tx},
		{&InternalError{Code: AuthFormEmailEmpty, Message: "Email is empty"}, AuthFormEmailEmpty},
	}

	for _, tc := range tests {
		code := tc.err.GetCode()
		assert.Equal(t, tc.expect, code, "Expected error code does not match actual")
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		code   int
		expect *InternalError
	}{
		{Tx, &InternalError{Code: Tx, Message: "[Transaction] Failed transaction"}},
		{TxCommit, &InternalError{Code: TxCommit, Message: "[Transaction] Failed to commit registration transaction"}},
		{TxNotUnique, &InternalError{Code: TxNotUnique, Message: "[Transaction] Value already exists"}},
		{TxUnknown, &InternalError{Code: TxUnknown, Message: "[Transaction] Failed transaction"}},
		{Auth, &InternalError{Code: Auth, Message: "[Auth] Error"}},
		{AuthDecode, &InternalError{Code: AuthDecode, Message: "[Auth]: Decode Error"}},
		{AuthForm, &InternalError{Code: AuthForm, Message: "[Auth]: Form data is invalid"}},
		{AuthFormEmailEmpty, &InternalError{Code: AuthFormEmailEmpty, Message: "[Auth]: Email is empty"}},
		{AuthFormEmailInvalid, &InternalError{Code: AuthFormEmailInvalid, Message: "[Auth]: Email address is invalid"}},
		{AuthFormPasswordInvalid, &InternalError{Code: AuthFormPasswordInvalid, Message: "[Auth]: Password is empty or less than 6 symbols"}},
		{AuthFormPasswordWrong, &InternalError{Code: AuthFormPasswordWrong, Message: "[Auth]: Wrong Password"}},
		{AuthFormPasswordsMismatch, &InternalError{Code: AuthFormPasswordsMismatch, Message: "[Auth]: Passwords mismatch"}},
		{QueryParamsToken, &InternalError{Code: QueryParamsToken, Message: "[Query Params]: Query parameter 'token' should be specified"}},
		{UserCredentials, &InternalError{Code: UserCredentials, Message: "[User Credentials]: Failed to get user credentials"}},
		{UserCredentialsToken, &InternalError{Code: UserCredentialsToken, Message: "[User Credentials]: User credentials with specified token does not exist"}},
		{UserCredentialsIsNotActive, &InternalError{Code: UserCredentialsIsNotActive, Message: "[User Credentials]: User is not activated"}},
		{UserCredentialsReset, &InternalError{Code: UserCredentialsReset, Message: "[User Credentials]: Failed to update user password"}},
		{Profile, &InternalError{Code: Profile, Message: "[Profile]: Failed to find user profile"}},
		{Token, &InternalError{Code: Token, Message: "[Token]: Token unknown error"}},
		{TokenEmpty, &InternalError{Code: TokenEmpty, Message: "[Token]: Token is empty"}},
		{TokenExpired, &InternalError{Code: TokenExpired, Message: "[Token]: Token expired, try to reset password"}},
		{JSON, &InternalError{Code: JSON, Message: "[JSON]: JSON unknown error"}},
		{JSONDecoder, &InternalError{Code: JSONDecoder, Message: "[JSON]: Decoder error"}},
		{DBGetUser, &InternalError{Code: DBGetUser, Message: "[DB]: Failed to get user"}},
		{Events, &InternalError{Code: Events, Message: "[Events]: Decode error"}},
		{EventsFightResult, &InternalError{Code: EventsFightResult, Message: "[Events]: Failed to set fight result"}},
		{EventIsDone, &InternalError{Code: EventIsDone, Message: "[Events]: Failed to set event done"}},
		{Count, &InternalError{Code: Count, Message: "[Count]: Failed to get items count"}},
		{CountFighters, &InternalError{Code: CountFighters, Message: "[Count]: Failed to get fighters count"}},
		{CountEvents, &InternalError{Code: CountEvents, Message: "[Count]: Failed to get events count"}},
		{Fighters, &InternalError{Code: Fighters, Message: "[Fighters]: Failed to find fighters"}},
		{Bets, &InternalError{Code: Bets, Message: "[Bets]: Error"}},
		{CountBets, &InternalError{Code: CountBets, Message: "[Bets]: Failed to get bets count"}},
		{9999, &InternalError{Code: 9999, Message: "Unknown error"}},
	}

	for _, tc := range tests {
		internalErr := New(tc.code)
		assert.Equal(t, tc.expect, internalErr, "Expected error does not match actual")
	}

}

func TestError(t *testing.T) {
	tests := []struct {
		err    *InternalError
		expect string
	}{
		{&InternalError{Code: Tx, Message: "Transaction error"}, "Transaction error"},
		{&InternalError{Code: AuthFormEmailEmpty, Message: "Email is empty"}, "Email is empty"},
	}

	for _, tc := range tests {
		message := tc.err.Error()
		assert.Equal(t, tc.expect, message, "Expected error message does not match actual")
	}
}
