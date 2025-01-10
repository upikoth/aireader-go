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

func (u *Users) Create(
	inputCtx context.Context,
	userToCreate *models.User,
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

	var dbCreatedUser ydbmodels.User
	dbUserToCreate := ydbmodels.NewYDBUserModel(userToCreate)

	err = u.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qRes, qErr := tx.QueryResultSet(
			qCtx,
			`declare $id as text;
			declare $email as text;
			declare $role as text;
			declare $password_hash as text;
			declare $created_at as timestamp;
			declare $updated_at as timestamp;
			declare $vk_id as text;
			declare $mailru_id as text;
			declare $yandex_id as text;

			insert into users
			(id, email, role, password_hash, created_at, updated_at, vk_id, mailru_id, yandex_id)
			values ($id, $email, $role, $password_hash, $created_at, $updated_at, $vk_id, $mailru_id, $yandex_id);

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
					Param("$id").Text(dbUserToCreate.ID).
					Param("$email").Text(dbUserToCreate.Email).
					Param("$role").Text(dbUserToCreate.Role).
					Param("$password_hash").Text(dbUserToCreate.PasswordHash).
					Param("$created_at").Timestamp(time.Now()).
					Param("$updated_at").Timestamp(time.Now()).
					Param("$vk_id").Text(dbUserToCreate.VkID).
					Param("$mailru_id").Text(dbUserToCreate.MailRuID).
					Param("$yandex_id").Text(dbUserToCreate.YandexID).
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
				query.Named("id", &dbCreatedUser.ID),
				query.Named("email", &dbCreatedUser.Email),
				query.Named("role", &dbCreatedUser.Role),
				query.Named("password_hash", &dbCreatedUser.PasswordHash),
				query.Named("vk_id", &dbCreatedUser.VkID),
				query.Named("mailru_id", &dbCreatedUser.MailRuID),
				query.Named("yandex_id", &dbCreatedUser.YandexID),
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

	return dbCreatedUser.FromYDBModel(), nil
}
