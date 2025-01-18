package voices

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

func (v *Voices) GetList(
	inputCtx context.Context,
	params *models.VoicesGetListParams,
) (res *models.VoiceList, err error) {
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

	var resVoices []*models.Voice
	var total int

	err = v.executeInQueryTransaction(ctx, func(qCtx context.Context, tx query.Transaction) error {
		qVoices, qErr := queryVoices(qCtx, tx, params)
		if qErr != nil {
			return qErr
		}
		resVoices = qVoices

		qTotal, qErr := queryVoicesTotal(qCtx, tx)
		if qErr != nil {
			return qErr
		}
		total = qTotal

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &models.VoiceList{
		Voices: resVoices,
		Total:  total,
	}, nil
}

func queryVoices(qCtx context.Context, tx query.Transaction, params *models.VoicesGetListParams) ([]*models.Voice, error) {
	var resVoices []*models.Voice

	qRes, qErr := tx.QueryResultSet(
		qCtx,
		`declare $limit as Uint64;
		declare $offset as Uint64;

		select
			name,
			source,
		from voices
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
		return resVoices, errors.WithStack(qErr)
	}

	defer func() { _ = qRes.Close(qCtx) }()

	for row, rErr := range qRes.Rows(qCtx) {
		if rErr != nil {
			return resVoices, errors.WithStack(rErr)
		}

		var voice ydbmodels.Voice

		sErr := row.ScanNamed(
			query.Named("name", &voice.Name),
			query.Named("source", &voice.Source),
		)

		if sErr != nil {
			return resVoices, errors.WithStack(sErr)
		}

		resVoices = append(resVoices, voice.FromYDBModel())
	}

	return resVoices, nil
}

func queryVoicesTotal(qCtx context.Context, tx query.Transaction) (int, error) {
	var total int
	qRes, qErr := tx.QueryResultSet(
		qCtx,
		`select count(*) as total from voices`,
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
