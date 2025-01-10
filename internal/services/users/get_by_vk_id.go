package users

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (u *Users) GetByVkID(
	inputCtx context.Context,
	vkID string,
) (*models.User, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	user, err := u.repositories.users.GetByVkID(ctx, vkID)

	if errors.Is(err, constants.ErrDBEntityNotFound) {
		return nil, constants.ErrUserNotFound
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return user, nil
}
