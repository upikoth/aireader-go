package ydbmodels

import "github.com/upikoth/aireader-go/internal/models"

type Registration struct {
	ID                string
	Email             string
	ConfirmationToken string
}

func NewYDBRegistrationModel(registration *models.Registration) *Registration {
	return &Registration{
		ID:                string(registration.ID),
		Email:             registration.Email,
		ConfirmationToken: registration.ConfirmationToken,
	}
}

func (r *Registration) FromYDBModel() *models.Registration {
	return &models.Registration{
		ID:                models.RegistrationID(r.ID),
		Email:             r.Email,
		ConfirmationToken: r.ConfirmationToken,
	}
}
