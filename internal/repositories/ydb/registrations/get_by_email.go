package registrations

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
)

func (r *Registrations) GetByEmail(
	inputCtx context.Context,
	email string,
) (res *models.Registration, err error) {
	return r.getBy(inputCtx, fieldNameGetByEmail, email)
}
