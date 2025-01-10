package users

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

func (u *Users) CreateByEmailPassword(
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

	createdUser, err := u.repositories.users.Create(
		ctx,
		newUser(
			email,
			withPasswordHash(string(passwordHash)),
		),
	)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return createdUser, nil
}
