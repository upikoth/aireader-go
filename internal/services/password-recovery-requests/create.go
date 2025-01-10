package passwordrecoveryrequests

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (p *PasswordRecoveryRequests) Create(
	inputCtx context.Context,
	params models.PasswordRecoveryRequestCreateParams,
) (*models.PasswordRecoveryRequest, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	passwordRecoveryRequest := newPasswordRecoveryRequest(params.Email)
	_, err := p.services.users.GetByEmail(ctx, passwordRecoveryRequest.Email)

	// Если пользователь не найден, возвращаем такой же ответ как если бы он был найден.
	// Так нельзя будет понять есть ли такой email в приложении.
	if errors.Is(err, constants.ErrUserNotFound) {
		return passwordRecoveryRequest, nil
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	existingPasswordRecoveryRequest, err :=
		p.repositories.passwordRecoveryRequests.GetByEmail(ctx, passwordRecoveryRequest.Email)

	if err != nil && !errors.Is(err, constants.ErrDBEntityNotFound) {
		tracing.HandleError(span, err)
		return nil, err
	}

	if err != nil && errors.Is(err, constants.ErrDBEntityNotFound) {
		passwordRecoveryRequest, err = p.repositories.passwordRecoveryRequests.Create(ctx, passwordRecoveryRequest)
	} else {
		passwordRecoveryRequest = existingPasswordRecoveryRequest
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	err = p.services.emails.SendPasswordRecoveryRequestEmail(
		ctx,
		passwordRecoveryRequest.Email,
		passwordRecoveryRequest.ConfirmationToken,
	)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	return passwordRecoveryRequest, nil
}
