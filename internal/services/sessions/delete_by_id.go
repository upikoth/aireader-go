package sessions

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (s *Sessions) DeleteByID(
	inputCtx context.Context,
	id models.SessionID,
) error {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	err := s.repositories.sessions.DeleteByID(ctx, id)

	if errors.Is(err, constants.ErrDBEntityNotFound) {
		return constants.ErrSessionNotFound
	}

	if err != nil {
		tracing.HandleError(span, err)
		return err
	}

	return nil
}
