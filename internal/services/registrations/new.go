package registrations

import (
	"github.com/google/uuid"
	"github.com/upikoth/aireader-go/internal/models"
)

type Option func(session *models.Registration)

func newRegistration(
	userEmail string,
	options ...Option,
) *models.Registration {
	registration := &models.Registration{
		ID:                models.RegistrationID(uuid.New().String()),
		Email:             userEmail,
		ConfirmationToken: uuid.New().String(),
	}

	for _, option := range options {
		option(registration)
	}

	return registration
}
