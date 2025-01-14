package voices

import (
	"context"
	"encoding/json"

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

func (v *Voices) Create(
	inputCtx context.Context,
	voiceToCreate *models.Voice,
) (res *models.Voice, err error) {
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

	var dbCreatedVoice ydbmodels.Voice
	dbVoiceToCreate := ydbmodels.NewYDBVoiceModel(voiceToCreate)

	err = v.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qRes, qErr := tx.QueryResultSet(
			qCtx,
			`declare $name as text;
			declare $source as text;

			insert into voices
			(name, source)
			values ($name, $source);

			select
				name,
				source,
			from voices
			where voices.name = $name;`,
			query.WithParameters(
				ydb.ParamsBuilder().
					Param("$name").Text(dbVoiceToCreate.Name).
					Param("$source").Text(dbVoiceToCreate.Source).
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
				query.Named("name", &dbCreatedVoice.Name),
				query.Named("source", &dbCreatedVoice.Source),
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

	return dbCreatedVoice.FromYDBModel(), nil
}
