package sessions

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (s *Sessions) CreateByUserID(
	inputCtx context.Context,
	userID models.UserID,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	session, err := s.repositories.sessions.Create(ctx, newSession(userID))

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return session, err
}
