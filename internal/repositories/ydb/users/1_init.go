package users

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type Users struct {
	db     *ydb.Driver
	qTx    query.Transaction
	logger logger.Logger
}

func New(
	db *ydb.Driver,
	logger logger.Logger,
) *Users {
	return &Users{
		db:     db,
		logger: logger,
	}
}

func (u *Users) WithTx(tx query.Transaction) *Users {
	return &Users{
		db:     u.db,
		qTx:    tx,
		logger: u.logger,
	}
}

func (u *Users) executeInQueryTransaction(
	ctx context.Context,
	fn func(c context.Context, tx query.Transaction) error,
) error {
	if u.qTx != nil {
		return fn(ctx, u.qTx)
	}

	return u.db.
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
