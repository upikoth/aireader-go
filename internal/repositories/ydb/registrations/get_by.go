package registrations

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/ydbmodels"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type fieldNameGetBy string

var (
	fieldNameGetByConfrimationToken fieldNameGetBy = "confirmation_token"
	fieldNameGetByEmail             fieldNameGetBy = "email"
)

func (r *Registrations) getBy(
	inputCtx context.Context,
	fieldName fieldNameGetBy,
	fieldValue string,
) (res *models.Registration, err error) {
	tracer := otel.Tracer(tracing.GetRepositoryYDBTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetRepositoryYDBTraceName())
	defer span.End()

	defer func() {
		if err != nil && !errors.Is(err, constants.ErrDBEntityNotFound) {
			span.RecordError(err)
			sentry.CaptureException(err)
		} else {
			bytes, _ := json.Marshal(res)
			span.SetAttributes(
				attribute.String("ydb.res", string(bytes)),
			)
		}
	}()

	var registration ydbmodels.Registration

	err = r.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qRes, qErr := tx.QueryResultSet(
			qCtx,
			fmt.Sprintf(
				`declare $filterValue as text;
				select
					id,
					email,
					confirmation_token,
				from registrations
				where %s = $filterValue;`,
				fieldName,
			),
			query.WithParameters(
				ydb.ParamsBuilder().Param("$filterValue").Text(fieldValue).Build(),
			),
		)

		if qErr != nil {
			return errors.WithStack(qErr)
		}

		defer func() { _ = qRes.Close(qCtx) }()

		for row, rErr := range qRes.Rows(qCtx) {
			if rErr != nil {
				return errors.WithStack(rErr)
			}

			sErr := row.ScanNamed(
				query.Named("id", &registration.ID),
				query.Named("email", &registration.Email),
				query.Named("confirmation_token", &registration.ConfirmationToken),
			)

			if sErr != nil {
				return errors.WithStack(sErr)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(registration).IsZero() {
		return nil, errors.WithStack(constants.ErrDBEntityNotFound)
	}

	return registration.FromYDBModel(), nil
}
