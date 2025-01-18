package voices

import (
	"context"

	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (v *Voices) GetList(
	inputCtx context.Context,
	params *models.VoicesGetListParams,
) (*models.VoiceList, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	res, err := v.repositories.voices.GetList(ctx, params)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return res, nil
}
