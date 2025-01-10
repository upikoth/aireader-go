package users

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (u *Users) CreateByEmailVkID(
	inputCtx context.Context,
	email string,
	vkID string,
) (*models.User, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	createdUser, err := u.repositories.users.Create(
		ctx,
		newUser(
			email,
			withVkID(vkID),
		),
	)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return createdUser, nil
}
