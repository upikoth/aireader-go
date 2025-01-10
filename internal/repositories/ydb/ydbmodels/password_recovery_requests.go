package ydbmodels

import "github.com/upikoth/aireader-go/internal/models"

type PasswordRecoveryRequest struct {
	ID                string
	Email             string
	ConfirmationToken string
}

func NewYDBPasswordRecoveryRequestModel(registration *models.PasswordRecoveryRequest) *PasswordRecoveryRequest {
	return &PasswordRecoveryRequest{
		ID:                string(registration.ID),
		Email:             registration.Email,
		ConfirmationToken: registration.ConfirmationToken,
	}
}

func (r *PasswordRecoveryRequest) FromYDBModel() *models.PasswordRecoveryRequest {
	return &models.PasswordRecoveryRequest{
		ID:                models.PasswordRecoveryRequestID(r.ID),
		Email:             r.Email,
		ConfirmationToken: r.ConfirmationToken,
	}
}
