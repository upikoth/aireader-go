package registrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"go.opentelemetry.io/otel"
)

func (r *Registrations) Confirm(
	inputCtx context.Context,
	params models.RegistrationConfirmParams,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	registration, err := r.repositories.registrations.GetByToken(ctx, params.ConfirmationToken)

	if errors.Is(err, constants.ErrDBEntityNotFound) {
		return nil, constants.ErrRegistrationNotFound
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	var createdUser *models.User
	err = r.Transaction(
		ctx,
		func(ctxTx context.Context, rTx *Registrations) error {
			dbErr := rTx.repositories.registrations.DeleteByID(ctxTx, registration.ID)

			if dbErr != nil {
				return dbErr
			}

			user, createUserErr := rTx.services.users.CreateByEmailPassword(
				ctxTx,
				registration.Email,
				params.Password,
			)
			createdUser = user

			if createUserErr != nil {
				return createUserErr
			}

			return nil
		},
		query.WithSerializableReadWrite(),
	)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	createdSession, err := r.services.sessions.CreateByUserID(
		ctx,
		createdUser.ID,
	)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, constants.ErrRegistrationCreatingSession
	}

	return createdSession, nil
}
