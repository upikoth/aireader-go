package sessions

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

func (s *Sessions) Create(
	inputCtx context.Context,
	sessionToCreate *models.Session,
) (res *models.SessionWithUserRole, err error) {
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

	var dbCreatedSession ydbmodels.Session
	dbSessionToCreate := ydbmodels.NewYDBSessionModel(sessionToCreate)

	err = s.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qRes, qErr := tx.QueryResultSet(
			qCtx,
			`declare $id as text;
			declare $token as text;
			declare $user_id as text;
			declare $created_at as timestamp;
			
			insert into sessions
			(id, token, user_id, created_at)
			values ($id, $token, $user_id, $created_at);

			select
				s.id as id,
				s.token as token,
				s.user_id as user_id,
				u.role as user_role,
			from sessions as s join users as u on s.user_id = u.id
			where s.id = $id;`,
			query.WithParameters(
				ydb.ParamsBuilder().
					Param("$id").Text(dbSessionToCreate.ID).
					Param("$token").Text(dbSessionToCreate.Token).
					Param("$user_id").Text(dbSessionToCreate.UserID).
					Param("$created_at").Timestamp(time.Now()).
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
				query.Named("id", &dbCreatedSession.ID),
				query.Named("token", &dbCreatedSession.Token),
				query.Named("user_id", &dbCreatedSession.UserID),
				query.Named("user_role", &dbCreatedSession.UserRole),
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

	return dbCreatedSession.FromYDBModel(), nil
}
