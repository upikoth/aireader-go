package users

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
	fieldNameGetByID       fieldNameGetBy = "id"
	fieldNameGetByEmail    fieldNameGetBy = "email"
	fieldNameGetByVkID     fieldNameGetBy = "vk_id"
	fieldNameGetByMailRuID fieldNameGetBy = "mailru_id"
	fieldNameGetByYandexID fieldNameGetBy = "yandex_id"
)

func (u *Users) getBy(
	inputCtx context.Context,
	fieldName fieldNameGetBy,
	fieldValue string,
) (res *models.User, err error) {
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

	var user ydbmodels.User

	err = u.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qRes, qErr := tx.QueryResultSet(
			qCtx,
			fmt.Sprintf(
				`declare $filterValue as text;
				select
					id,
					email,
					role,
					password_hash,
					vk_id,
					mailru_id,
					yandex_id,
				from users
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
				query.Named("id", &user.ID),
				query.Named("email", &user.Email),
				query.Named("role", &user.Role),
				query.Named("password_hash", &user.PasswordHash),
				query.Named("vk_id", &user.VkID),
				query.Named("mailru_id", &user.MailRuID),
				query.Named("yandex_id", &user.YandexID),
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

	if reflect.ValueOf(user).IsZero() {
		return nil, errors.WithStack(constants.ErrDBEntityNotFound)
	}

	return user.FromYDBModel(), nil
}
