package users

import (
	"context"
	"encoding/json"
	"strconv"

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

func (u *Users) GetList(
	inputCtx context.Context,
	params *models.UsersGetListParams,
) (res *models.UserList, err error) {
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

	var resUsers []*models.User
	var total int

	err = u.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qUsers, qErr := queryUsers(qCtx, tx, params)
		if qErr != nil {
			return qErr
		}
		resUsers = qUsers

		qTotal, qErr := queryUsersTotal(qCtx, tx)
		if qErr != nil {
			return qErr
		}
		total = qTotal

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &models.UserList{
		Users: resUsers,
		Total: total,
	}, nil
}

func queryUsers(qCtx context.Context, tx query.Transaction, params *models.UsersGetListParams) ([]*models.User, error) {
	var resUsers []*models.User

	qRes, qErr := tx.QueryResultSet(
		qCtx,
		`declare $limit as Uint64;
		declare $offset as Uint64;

		select
			id,
			email,
			password_hash,
			role,
			vk_id,
			mailru_id,
			yandex_id,
		from users
		limit $limit
		offset $offset`,
		query.WithParameters(
			ydb.ParamsBuilder().
				Param("$limit").Uint64(uint64(params.Limit)).
				Param("$offset").Uint64(uint64(params.Offset)).
				Build(),
		),
	)

	if qErr != nil {
		return resUsers, errors.WithStack(qErr)
	}

	defer func() { _ = qRes.Close(qCtx) }()

	for row, rErr := range qRes.Rows(qCtx) {
		if rErr != nil {
			return resUsers, errors.WithStack(rErr)
		}

		var user ydbmodels.User

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
			return resUsers, errors.WithStack(sErr)
		}

		resUsers = append(resUsers, user.FromYDBModel())
	}

	return resUsers, nil
}

func queryUsersTotal(qCtx context.Context, tx query.Transaction) (int, error) {
	var total int
	qRes, qErr := tx.QueryResultSet(
		qCtx,
		`select count(*) as total from users`,
	)

	if qErr != nil {
		return total, errors.WithStack(qErr)
	}

	for row, rErr := range qRes.Rows(qCtx) {
		if rErr != nil {
			return total, errors.WithStack(rErr)
		}

		var qTotal uint64
		sErr := row.ScanNamed(
			query.Named("total", &qTotal),
		)

		if sErr != nil {
			return total, errors.WithStack(sErr)
		}

		total, _ = strconv.Atoi(strconv.FormatUint(qTotal, 10))
	}

	return total, nil
}
