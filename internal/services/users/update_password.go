package users

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

func (u *Users) UpdatePassword(
	inputCtx context.Context,
	email string,
	password string,
) (*models.User, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	user, err := u.GetByEmail(ctx, email)

	if errors.Is(err, constants.ErrUserNotFound) {
		return nil, constants.ErrUserNotFound
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	user.PasswordHash = string(passwordHash)

	updatedUser, err := u.repositories.users.Update(ctx, user)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return updatedUser, nil
}
