package voices

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (v *Voices) GetByName(
	inputCtx context.Context,
	name string,
) (*models.Voice, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	voice, err := v.repositories.voices.GetByName(ctx, name)

	if errors.Is(err, constants.ErrDBEntityNotFound) {
		return nil, constants.ErrVoiceNotFound
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return voice, nil
}
