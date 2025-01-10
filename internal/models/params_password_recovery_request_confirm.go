package models

type PasswordRecoveryRequestConfirmParams struct {
	ConfirmationToken string
	NewPassword       string
}
