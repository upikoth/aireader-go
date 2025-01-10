package passwordrecoveryrequests

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"go.opentelemetry.io/otel"
)

func (p *PasswordRecoveryRequests) Confirm(
	inputCtx context.Context,
	params models.PasswordRecoveryRequestConfirmParams,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	passwordRecoveryRequest, err := p.
		repositories.
		passwordRecoveryRequests.
		GetByToken(ctx, params.ConfirmationToken)

	if errors.Is(err, constants.ErrDBEntityNotFound) {
		return nil, constants.ErrPasswordRecoveryRequestNotFound
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	var updatedUser *models.User
	err = p.Transaction(
		ctx,
		func(ctxTx context.Context, pTx *PasswordRecoveryRequests) error {
			user, updateErr := pTx.services.users.UpdatePassword(ctxTx, passwordRecoveryRequest.Email, params.NewPassword)
			updatedUser = user

			if updateErr != nil {
				return updateErr
			}

			deleteErr := pTx.repositories.passwordRecoveryRequests.DeleteByID(
				ctxTx,
				passwordRecoveryRequest.ID,
			)

			if deleteErr != nil {
				return deleteErr
			}

			return nil
		},
		query.WithSerializableReadWrite(),
	)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	createdSession, err := p.services.sessions.CreateByUserID(ctx, updatedUser.ID)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, constants.ErrPasswordRecoveryRequestCreatingSession
	}

	return createdSession, nil
}
