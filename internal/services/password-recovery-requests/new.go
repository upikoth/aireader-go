package passwordrecoveryrequests

import (
	"github.com/google/uuid"
	"github.com/upikoth/aireader-go/internal/models"
)

type Option func(session *models.PasswordRecoveryRequest)

func newPasswordRecoveryRequest(
	email string,
	options ...Option,
) *models.PasswordRecoveryRequest {
	prr := &models.PasswordRecoveryRequest{
		ID:                models.PasswordRecoveryRequestID(uuid.New().String()),
		Email:             email,
		ConfirmationToken: uuid.New().String(),
	}

	for _, option := range options {
		option(prr)
	}

	return prr
}
