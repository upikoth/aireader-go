package voices

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
)

func (v *Voices) GetByName(
	inputCtx context.Context,
	name string,
) (res *models.Voice, err error) {
	return v.getBy(inputCtx, fieldNameGetByName, name)
}
