package users

import (
	"context"
	"encoding/json"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/ydbmodels"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

//nolint:funlen
func (u *Users) Update(
	inputCtx context.Context,
	userToUpdate *models.User,
) (res *models.User, err error) {
	tracer := otel.Tracer(tracing.GetRepositoryYDBTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetRepositoryYDBTraceName())
	defer span.End()

	defer func() {
		if err != nil {
			span.RecordError(err)
			sentry.CaptureException(err)
		} else {
			bytes, _ := json.Marshal(res)
			span.SetAttributes(
				attribute.String("ydb.res", string(bytes)),
			)
		}
	}()

	if userToUpdate == nil || userToUpdate.ID == "" {
		return nil, errors.New("не задан id пользователя")
	}

	var dbUpdatedUser ydbmodels.User
	dbUserToUpdate := ydbmodels.NewYDBUserModel(userToUpdate)

	err = u.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qRes, qErr := tx.QueryResultSet(
			qCtx,
			`declare $id as text;
			declare $email as text;
			declare $role as text;
			declare $password_hash as text;
			declare $updated_at as timestamp;
			declare $vk_id as text;
			declare $mailru_id as text;
			declare $yandex_id as text;

			update users
			set
				email = $email,
				role = $role,
				password_hash = $password_hash,
				updated_at = $updated_at,
				vk_id = $vk_id,
				mailru_id = $mailru_id,
				yandex_id = $yandex_id
			where id = $id;

			select
				id,
				email,
				role,
				password_hash,
				vk_id,
				mailru_id,
				yandex_id,
			from users
			where users.id = $id;`,
			query.WithParameters(
				ydb.ParamsBuilder().
					Param("$id").Text(dbUserToUpdate.ID).
					Param("$email").Text(dbUserToUpdate.Email).
					Param("$role").Text(dbUserToUpdate.Role).
					Param("$password_hash").Text(dbUserToUpdate.PasswordHash).
					Param("$updated_at").Timestamp(time.Now()).
					Param("$vk_id").Text(dbUserToUpdate.VkID).
					Param("$mailru_id").Text(dbUserToUpdate.MailRuID).
					Param("$yandex_id").Text(dbUserToUpdate.YandexID).
					Build(),
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
				query.Named("id", &dbUpdatedUser.ID),
				query.Named("email", &dbUpdatedUser.Email),
				query.Named("role", &dbUpdatedUser.Role),
				query.Named("password_hash", &dbUpdatedUser.PasswordHash),
				query.Named("vk_id", &dbUpdatedUser.VkID),
				query.Named("mailru_id", &dbUpdatedUser.MailRuID),
				query.Named("yandex_id", &dbUpdatedUser.YandexID),
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

	return dbUpdatedUser.FromYDBModel(), nil
}
