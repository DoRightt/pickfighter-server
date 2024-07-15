package internal

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name   string
		err    *Error
		expect string
	}{
		{
			"Default message",
			&Error{Code: Tx, Message: "Just Error"},
			"Just Error",
		},
		{
			"Empty message",
			&Error{Code: AuthFormEmailEmpty, Message: ""},
			"",
		},
	}

	for _, tc := range tests {
		msg := tc.err.GetMessage()
		assert.Equal(t, tc.expect, msg, "Expected error code does not match actual")
	}
}

func TestGetCode(t *testing.T) {
	tests := []struct {
		err    *Error
		expect int
	}{
		{
			&Error{Code: Tx, Message: "Transaction error"},
			Tx,
		},
		{
			&Error{Code: AuthFormEmailEmpty, Message: "Email is empty"},
			AuthFormEmailEmpty,
		},
	}

	for _, tc := range tests {
		code := tc.err.GetCode()
		assert.Equal(t, tc.expect, code, "Expected error code does not match actual")
	}
}

func TestNewDefault(t *testing.T) {
	tests := []struct {
		code   int
		expect *Error
	}{
		{Tx, &Error{Code: Tx, Message: "[Transaction] Failed transaction"}},
		{TxCommit, &Error{Code: TxCommit, Message: "[Transaction] Failed to commit registration transaction"}},
		{TxNotUnique, &Error{Code: TxNotUnique, Message: "[Transaction] Value already exists"}},
		{TxUnknown, &Error{Code: TxUnknown, Message: "[Transaction] Failed transaction"}},
		{Auth, &Error{Code: Auth, Message: "[Auth] Error"}},
		{AuthDecode, &Error{Code: AuthDecode, Message: "[Auth]: Decode Error"}},
		{AuthForm, &Error{Code: AuthForm, Message: "[Auth]: Form data is invalid"}},
		{AuthFormEmailEmpty, &Error{Code: AuthFormEmailEmpty, Message: "[Auth]: Email is empty"}},
		{AuthFormEmailInvalid, &Error{Code: AuthFormEmailInvalid, Message: "[Auth]: Email address is invalid"}},
		{AuthFormPasswordInvalid, &Error{Code: AuthFormPasswordInvalid, Message: "[Auth]: Password is empty or less than 6 symbols"}},
		{AuthFormPasswordWrong, &Error{Code: AuthFormPasswordWrong, Message: "[Auth]: Wrong Password"}},
		{AuthFormPasswordsMismatch, &Error{Code: AuthFormPasswordsMismatch, Message: "[Auth]: Passwords mismatch"}},
		{QueryParamsToken, &Error{Code: QueryParamsToken, Message: "[Query Params]: Query parameter 'token' should be specified"}},
		{UserCredentials, &Error{Code: UserCredentials, Message: "[User Credentials]: Failed to get user credentials"}},
		{UserCredentialsToken, &Error{Code: UserCredentialsToken, Message: "[User Credentials]: User credentials with specified token does not exist"}},
		{UserCredentialsIsNotActive, &Error{Code: UserCredentialsIsNotActive, Message: "[User Credentials]: User is not activated"}},
		{UserCredentialsReset, &Error{Code: UserCredentialsReset, Message: "[User Credentials]: Failed to update user password"}},
		{Profile, &Error{Code: Profile, Message: "[Profile]: Failed to find user profile"}},
		{Token, &Error{Code: Token, Message: "[Token]: Token unknown error"}},
		{TokenEmpty, &Error{Code: TokenEmpty, Message: "[Token]: Token is empty"}},
		{TokenExpired, &Error{Code: TokenExpired, Message: "[Token]: Token expired, try to reset password"}},
		{JSON, &Error{Code: JSON, Message: "[JSON]: JSON unknown error"}},
		{JSONDecoder, &Error{Code: JSONDecoder, Message: "[JSON]: Decoder error"}},
		{DBGetUser, &Error{Code: DBGetUser, Message: "[DB]: Failed to get user"}},
		{Count, &Error{Code: Count, Message: "[Count]: Failed to get items count"}},
		{CountFighters, &Error{Code: CountFighters, Message: "[Count]: Failed to get fighters count"}},
		{Fighters, &Error{Code: Fighters, Message: "[Fighters]: Failed to find fighters"}},
		{9999, &Error{Code: 9999, Message: "Unknown Error"}},
	}

	for _, tc := range tests {
		internalErr := NewDefault(tc.code, 0)
		assert.Equal(t, tc.expect.Message, internalErr.Message, "Expected error does not match actual")
	}

}

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		err          error
		internal     int
		expectedCode int
		expectedMsg  string
		expectedInt  int
	}{
		{
			name:         "Valid error",
			code:         404,
			err:          errors.New("Not Found"),
			internal:     123,
			expectedCode: 404,
			expectedMsg:  "Not Found",
			expectedInt:  123,
		},
		{
			name:         "Another valid error",
			code:         500,
			err:          errors.New("Internal Server Error"),
			internal:     456,
			expectedCode: 500,
			expectedMsg:  "Internal Server Error",
			expectedInt:  456,
		},
		{
			name:         "Error with empty message",
			code:         400,
			err:          errors.New(""),
			internal:     789,
			expectedCode: 400,
			expectedMsg:  "",
			expectedInt:  789,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := New(tc.code, tc.err, tc.internal)

			assert.Equal(t, tc.expectedCode, e.Code)
			assert.Equal(t, tc.expectedMsg, e.Message)
			assert.Equal(t, tc.expectedInt, e.InternalCode)

			timestampStr := e.Timestamp

			_, err := time.Parse(time.RFC1123, timestampStr)
			assert.NoError(t, err, "Timestamp format is incorrect")
		})
	}
}

func TestError(t *testing.T) {
	// Define the table of test cases
	tests := []struct {
		name          string
		code          int
		err           error
		internal      int
		timestamp     string
		expectedError string
	}{
		{
			name:          "Valid error",
			code:          404,
			err:           errors.New("Not Found"),
			internal:      123,
			timestamp:     "Mon, 01 Jan 2000 00:00:00 UTC",
			expectedError: "[ERROR]: Not Found. \n[ERROR CODE]: 404. \n[INTERNAL CODE]: 123. \nSERVICE: fighters-service.\nTime: Mon, 01 Jan 2000 00:00:00 UTC.\n",
		},
		{
			name:          "Another valid error",
			code:          500,
			err:           errors.New("Internal Server Error"),
			internal:      456,
			timestamp:     "Tue, 02 Feb 2001 12:34:56 UTC",
			expectedError: "[ERROR]: Internal Server Error. \n[ERROR CODE]: 500. \n[INTERNAL CODE]: 456. \nSERVICE: fighters-service.\nTime: Tue, 02 Feb 2001 12:34:56 UTC.\n",
		},
		{
			name:          "Error with empty message",
			code:          400,
			err:           errors.New(""),
			internal:      789,
			timestamp:     "Wed, 03 Mar 2002 23:45:01 UTC",
			expectedError: "[ERROR]: . \n[ERROR CODE]: 400. \n[INTERNAL CODE]: 789. \nSERVICE: fighters-service.\nTime: Wed, 03 Mar 2002 23:45:01 UTC.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the Error object with the current test case inputs
			e := &Error{
				Code:         tt.code,
				InternalCode: tt.internal,
				Message:      tt.err.Error(),
				Timestamp:    tt.timestamp,
			}

			// Check if the returned error message matches expected value
			assert.Equal(t, tt.expectedError, e.Error())
		})
	}
}
