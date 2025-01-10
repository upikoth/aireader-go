package registrations

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
)

func (r *Registrations) GetByToken(
	inputCtx context.Context,
	confirmationToken string,
) (res *models.Registration, err error) {
	return r.getBy(inputCtx, fieldNameGetByConfrimationToken, confirmationToken)
}
