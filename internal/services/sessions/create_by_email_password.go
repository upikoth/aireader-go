package sessions

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

func (s *Sessions) CreateByEmailPassword(
	inputCtx context.Context,
	email string,
	password string,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	user, err := s.services.users.GetByEmail(ctx, email)

	if errors.Is(err, constants.ErrUserNotFound) {
		return nil, constants.ErrSessionCreateInvalidCredentials
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return nil, constants.ErrSessionCreateInvalidCredentials
	}

	session, err := s.repositories.sessions.Create(ctx, newSession(user.ID))

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return session, err
}
