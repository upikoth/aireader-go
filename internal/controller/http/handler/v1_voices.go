//nolint:dupl // тут нужно дублировать
package handler

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	app "github.com/upikoth/aireader-go/internal/generated/app"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (h *Handler) V1CreateVoice(
	inputCtx context.Context,
	req *app.V1VoicesCreateVoiceRequestBody,
	params app.V1CreateVoiceParams,
) (*app.SuccessResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	session, err := h.services.Sessions.CheckToken(ctx, params.AuthorizationToken)

	if errors.Is(err, constants.ErrSessionNotFound) {
		return nil, &models.Error{
			Code:        models.ErrCodeUserUnauthorized,
			Description: "User session is invalid",
			StatusCode:  http.StatusUnauthorized,
		}
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	if !session.UserRole.CheckAccessToAction(models.UserActionVoiceCreate) {
		return nil, &models.Error{
			Code:        models.ErrCodeVoiceCreateForbidden,
			Description: "Insufficient rights",
			StatusCode:  http.StatusForbidden,
		}
	}

	voiceCreateParams := models.VoiceCreateParams{
		Name:   req.Name,
		Source: models.VoiceSource(req.Source),
	}

	_, err = h.services.Voices.Create(ctx, voiceCreateParams)

	if errors.Is(err, constants.ErrVoiceNameAlreadyExist) {
		return nil, &models.Error{
			Code:        models.ErrCodeVoiceAlreadyExist,
			Description: "A voice with that name already exists",
			StatusCode:  http.StatusBadRequest,
		}
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.SuccessResponse{
		Success: true,
		Data:    app.SuccessResponseData{},
	}, nil
}

func (h *Handler) V1GetVoices(
	inputCtx context.Context,
	params app.V1GetVoicesParams,
) (*app.V1VoicesGetVoicesResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	voicesGetListParams := &models.VoicesGetListParams{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	voiceList, err := h.services.Voices.GetList(ctx, voicesGetListParams)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	var voicesResult []app.Voice
	for _, voice := range voiceList.Voices {
		voicesResult = append(voicesResult, app.Voice{
			Name:   voice.Name,
			Source: app.VoiceSource(voice.Source),
		})
	}

	return &app.V1VoicesGetVoicesResponse{
		Success: true,
		Data: app.V1VoicesGetVoicesResponseData{
			Voices: voicesResult,
			Limit:  voicesGetListParams.Limit,
			Offset: voicesGetListParams.Offset,
			Total:  voiceList.Total,
		},
	}, nil
}
