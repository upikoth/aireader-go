package users

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (u *Users) UpdateUser(
	inputCtx context.Context,
	user *models.User,
) (*models.User, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	updatedUser, err := u.repositories.users.Update(ctx, user)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return updatedUser, nil
}
