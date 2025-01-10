package models

type RegistrationID string

type Registration struct {
	ID                RegistrationID
	Email             string
	ConfirmationToken string
}
