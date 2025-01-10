package passwordrecoveryrequests

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
)

func (p *PasswordRecoveryRequests) GetByEmail(
	inputCtx context.Context,
	email string,
) (res *models.PasswordRecoveryRequest, err error) {
	return p.getBy(inputCtx, fieldNameGetByEmail, email)
}
