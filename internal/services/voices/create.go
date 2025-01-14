package voices

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (v *Voices) Create(
	inputCtx context.Context,
	params models.VoiceCreateParams,
) (*models.Voice, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	_, err := v.GetByName(ctx, params.Name)

	// Если есть ошибка, которая отличается от того что голос не найден.
	if err != nil && !errors.Is(err, constants.ErrVoiceNotFound) {
		tracing.HandleError(span, err)
		return nil, err
	}

	// Если голос найден.
	if err == nil {
		return nil, constants.ErrVoiceNameAlreadyExist
	}

	voice, err := v.repositories.voices.Create(ctx, newVoice(params.Name, params.Source))

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return voice, nil
}
