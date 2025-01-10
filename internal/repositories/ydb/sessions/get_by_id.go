package sessions

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
)

func (s *Sessions) GetByID(
	inputCtx context.Context,
	id string,
) (res *models.SessionWithUserRole, err error) {
	return s.getBy(inputCtx, fieldNameGetByID, id)
}
