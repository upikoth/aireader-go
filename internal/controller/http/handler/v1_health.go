package handler

import (
	"context"

	app "github.com/upikoth/aireader-go/internal/generated/app"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (h *Handler) V1CheckHealth(inputCtx context.Context) (*app.SuccessResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	_, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	return &app.SuccessResponse{
		Success: app.SuccessResponseSuccessTrue,
	}, nil
}
