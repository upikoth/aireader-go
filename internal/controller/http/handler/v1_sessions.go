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

func (h *Handler) V1CheckCurrentSession(
	inputCtx context.Context,
	params app.V1CheckCurrentSessionParams,
) (*app.SuccessResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	_, err := h.services.Sessions.CheckToken(ctx, params.AuthorizationToken)

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

	return &app.SuccessResponse{
		Success: app.SuccessResponseSuccessTrue,
	}, nil
}

func (h *Handler) V1CreateSession(
	inputCtx context.Context,
	req *app.V1SessionsCreateSessionRequestBody,
) (*app.V1SessionsCreateSessionResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	session, err := h.services.Sessions.CreateByEmailPassword(ctx, req.Email, string(req.Password))

	if errors.Is(err, constants.ErrSessionCreateInvalidCredentials) {
		return nil, &models.Error{
			Code:        models.ErrCodeCreateSessionWrongEmailOrPassword,
			Description: "Incorrect email or password",
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

	return &app.V1SessionsCreateSessionResponse{
		Success: true,
		Data: app.V1SessionsCreateSessionResponseData{
			Session: app.Session{
				ID:       string(session.ID),
				Token:    session.Token,
				UserRole: app.UserRole(session.UserRole),
			},
		},
	}, nil
}

func (h *Handler) V1DeleteSession(
	inputCtx context.Context,
	params app.V1DeleteSessionParams,
) (*app.SuccessResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	err := h.services.Sessions.DeleteByID(ctx, models.SessionID(params.ID))

	if errors.Is(err, constants.ErrSessionNotFound) {
		return nil, &models.Error{
			Code:        models.ErrCodeSessionNotFound,
			Description: "Session with the given id was not found",
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
		Success: app.SuccessResponseSuccessTrue,
	}, nil
}
