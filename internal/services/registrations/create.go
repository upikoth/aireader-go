package registrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (r *Registrations) Create(
	inputCtx context.Context,
	params models.RegistrationCreateParams,
) (*models.Registration, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	_, err := r.services.users.GetByEmail(ctx, params.Email)

	// Если есть ошибка, которая отличается от того что пользователь не найден.
	if err != nil && !errors.Is(err, constants.ErrUserNotFound) {
		tracing.HandleError(span, err)
		return nil, err
	}

	// Если пользователь найден.
	if err == nil {
		return nil, constants.ErrUserAlreadyExist
	}

	existingRegistration, err := r.repositories.registrations.GetByEmail(ctx, params.Email)

	// Если есть ошибка, которая отличается от того что регистрация не найдена.
	if err != nil && !errors.Is(err, constants.ErrDBEntityNotFound) {
		tracing.HandleError(span, err)
		return nil, err
	}

	var registration *models.Registration

	if err != nil && errors.Is(err, constants.ErrDBEntityNotFound) {
		registration, err = r.repositories.registrations.Create(ctx, newRegistration(params.Email))
	} else {
		registration = existingRegistration
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	err = r.services.emails.SendRegistrationEmail(
		ctx,
		registration.Email,
		registration.ConfirmationToken,
	)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return registration, nil
}
