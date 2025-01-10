package models

type PasswordRecoveryRequestID string

type PasswordRecoveryRequest struct {
	ID                PasswordRecoveryRequestID
	Email             string
	ConfirmationToken string
}
