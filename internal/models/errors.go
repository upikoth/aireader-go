package models

import "net/http"

type ErrorCode string

const (
	ErrCodeValidationByOpenapi ErrorCode = "1"
	ErrCodeUserUnauthorized    ErrorCode = "2"
	ErrCodeInterval            ErrorCode = "3"

	ErrCodeRegistrationUserWithThisEmailAlreadyExist ErrorCode = "100"
	ErrCodeRegistrationNotFound                      ErrorCode = "101"
	ErrCodeRegistrationCreatingSession               ErrorCode = "102"

	ErrCodeCreateSessionWrongEmailOrPassword ErrorCode = "200"
	ErrCodeSessionNotFound                   ErrorCode = "201"

	ErrCodePasswordRecoveryRequestNotFound        ErrorCode = "300"
	ErrCodePasswordRecoveryRequestCreatingSession ErrorCode = "301"

	ErrCodeUsersGetListForbidden ErrorCode = "400"

	ErrCodeOauthSourceNotExist ErrorCode = "500"
)

type Error struct {
	Code        ErrorCode
	StatusCode  int
	Description string
}

func (e *Error) Error() string {
	return e.Description
}

func (e *Error) GetStatusCode() int {
	if e.StatusCode == 0 {
		return http.StatusInternalServerError
	}

	return e.StatusCode
}
