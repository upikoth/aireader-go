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

func (h *Handler) V1CreateRegistration(
	inputCtx context.Context,
	req *app.V1RegistrationsCreateRegistrationRequestBody,
) (*app.V1RegistrationsCreateRegistrationResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	registrationCreateParams := models.RegistrationCreateParams{
		Email: req.Email,
	}

	registration, err := h.services.Registrations.Create(ctx, registrationCreateParams)

	if errors.Is(err, constants.ErrUserAlreadyExist) {
		return nil, &models.Error{
			Code:        models.ErrCodeRegistrationUserWithThisEmailAlreadyExist,
			Description: "A user with the specified email already exists",
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

	return &app.V1RegistrationsCreateRegistrationResponse{
		Success: true,
		Data: app.V1RegistrationsCreateRegistrationResponseData{
			ID:    string(registration.ID),
			Email: registration.Email,
		},
	}, nil
}

func (h *Handler) V1ConfirmRegistration(
	inputCtx context.Context,
	req *app.V1RegistrationsConfirmRegistrationRequestBody,
) (*app.V1RegistrationsConfirmRegistrationResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	registrationConfirmParams := models.RegistrationConfirmParams{
		ConfirmationToken: req.ConfirmationToken,
		Password:          string(req.Password),
	}

	session, err := h.services.Registrations.Confirm(ctx, registrationConfirmParams)

	if errors.Is(err, constants.ErrRegistrationNotFound) {
		return nil, &models.Error{
			Code:        models.ErrCodeRegistrationNotFound,
			Description: "Registration with transferred token not found",
			StatusCode:  http.StatusBadRequest,
		}
	}

	if errors.Is(err, constants.ErrRegistrationCreatingSession) {
		return nil, &models.Error{
			Code:        models.ErrCodeRegistrationCreatingSession,
			Description: "Session not created",
		}
	}

	if err != nil {
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.V1RegistrationsConfirmRegistrationResponse{
		Success: true,
		Data: app.V1RegistrationsConfirmRegistrationResponseData{
			Session: app.Session{
				ID:       string(session.ID),
				Token:    session.Token,
				UserRole: app.UserRole(session.UserRole),
			},
		},
	}, nil
}
