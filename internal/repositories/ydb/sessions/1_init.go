package sessions

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type Sessions struct {
	db     *ydb.Driver
	qTx    query.Transaction
	logger logger.Logger
}

func New(
	db *ydb.Driver,
	logger logger.Logger,
) *Sessions {
	return &Sessions{
		db:     db,
		logger: logger,
	}
}

func (s *Sessions) WithTx(tx query.Transaction) *Sessions {
	return &Sessions{
		db:     s.db,
		qTx:    tx,
		logger: s.logger,
	}
}

func (s *Sessions) executeInQueryTransaction(
	ctx context.Context,
	fn func(c context.Context, tx query.Transaction) error,
) error {
	if s.qTx != nil {
		return fn(ctx, s.qTx)
	}

	return s.db.
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
