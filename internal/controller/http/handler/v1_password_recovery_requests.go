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

func (h *Handler) V1CreatePasswordRecoveryRequest(
	inputCtx context.Context,
	req *app.V1PasswordRecoveryRequestsCreatePasswordRecoveryRequestRequestBody,
) (*app.V1PasswordRecoveryRequestsCreatePasswordRecoveryRequestResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	passwordRecoveryRequestCreateParams := models.PasswordRecoveryRequestCreateParams{
		Email: req.Email,
	}

	passwordRecoveryRequest, err := h.services.PasswordRecoveryRequests.Create(ctx, passwordRecoveryRequestCreateParams)

	if err != nil {
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.V1PasswordRecoveryRequestsCreatePasswordRecoveryRequestResponse{
		Success: true,
		Data: app.V1PasswordRecoveryRequestsCreatePasswordRecoveryRequestResponseData{
			ID:    string(passwordRecoveryRequest.ID),
			Email: passwordRecoveryRequest.Email,
		},
	}, nil
}

func (h *Handler) V1ConfirmPasswordRecoveryRequest(
	inputCtx context.Context,
	req *app.V1PasswordRecoveryRequestsConfirmPasswordRecoveryRequestRequestBody,
) (*app.V1PasswordRecoveryRequestsConfirmPasswordRecoveryRequestResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	passwordRecoveryRequestConfirmParams := models.PasswordRecoveryRequestConfirmParams{
		ConfirmationToken: req.ConfirmationToken,
		NewPassword:       string(req.NewPassword),
	}

	session, err := h.services.PasswordRecoveryRequests.Confirm(ctx, passwordRecoveryRequestConfirmParams)

	if errors.Is(err, constants.ErrPasswordRecoveryRequestNotFound) {
		return nil, &models.Error{
			Code:        models.ErrCodePasswordRecoveryRequestNotFound,
			Description: "Password recovery request with transferred token not found",
			StatusCode:  http.StatusBadRequest,
		}
	}

	if errors.Is(err, constants.ErrPasswordRecoveryRequestCreatingSession) {
		return nil, &models.Error{
			Code:        models.ErrCodePasswordRecoveryRequestCreatingSession,
			Description: "Session not created",
		}
	}

	if err != nil {
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.V1PasswordRecoveryRequestsConfirmPasswordRecoveryRequestResponse{
		Success: true,
		Data: app.V1PasswordRecoveryRequestsConfirmPasswordRecoveryRequestResponseData{
			Session: app.Session{
				ID:       string(session.ID),
				Token:    session.Token,
				UserRole: app.UserRole(session.UserRole),
			},
		},
	}, nil
}
