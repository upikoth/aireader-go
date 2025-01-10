package users

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
)

func (u *Users) GetByID(
	inputCtx context.Context,
	id models.UserID,
) (res *models.User, err error) {
	return u.getBy(inputCtx, fieldNameGetByID, string(id))
}
