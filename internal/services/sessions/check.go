package sessions

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (s *Sessions) CheckToken(
	inputCtx context.Context,
	token string,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	session, err := s.repositories.sessions.GetByToken(ctx, token)

	if errors.Is(err, constants.ErrDBEntityNotFound) {
		return nil, constants.ErrSessionNotFound
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return session, err
}
