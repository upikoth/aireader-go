package voices

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type Voices struct {
	db     *ydb.Driver
	qTx    query.Transaction
	logger logger.Logger
}

func New(
	db *ydb.Driver,
	logger logger.Logger,
) *Voices {
	return &Voices{
		db:     db,
		logger: logger,
	}
}

func (v *Voices) WithTx(tx query.Transaction) *Voices {
	return &Voices{
		db:     v.db,
		qTx:    tx,
		logger: v.logger,
	}
}

func (v *Voices) executeInQueryTransaction(
	ctx context.Context,
	fn func(c context.Context, tx query.Transaction) error,
) error {
	if v.qTx != nil {
		return fn(ctx, v.qTx)
	}

	return v.db.
		Query().
		Do(ctx, func(qCtx context.Context, s query.Session) error {
			qTx, qErr := s.Begin(qCtx, query.TxSettings(query.WithSerializableReadWrite()))

			if qErr != nil {
				return errors.WithStack(qErr)
			}
			defer func() { _ = qTx.Rollback(qCtx) }()

			if fnErr := fn(qCtx, qTx); fnErr != nil {
				return fnErr
			}

			return qTx.CommitTx(qCtx)
		})
}
