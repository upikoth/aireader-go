package users

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
)

func (u *Users) GetByYandexID(
	inputCtx context.Context,
	yandexID string,
) (res *models.User, err error) {
	return u.getBy(inputCtx, fieldNameGetByYandexID, yandexID)
}
