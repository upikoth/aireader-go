package registrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type Registrations struct {
	db     *ydb.Driver
	qTx    query.Transaction
	logger logger.Logger
}

func New(
	db *ydb.Driver,
	logger logger.Logger,
) *Registrations {
	return &Registrations{
		db:     db,
		logger: logger,
	}
}

func (r *Registrations) WithTx(tx query.Transaction) *Registrations {
	return &Registrations{
		db:     r.db,
		qTx:    tx,
		logger: r.logger,
	}
}

func (r *Registrations) executeInQueryTransaction(
	ctx context.Context,
	fn func(c context.Context, tx query.Transaction) error,
) error {
	if r.qTx != nil {
		return fn(ctx, r.qTx)
	}

	return r.db.
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
