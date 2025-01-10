package passwordrecoveryrequests

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type PasswordRecoveryRequests struct {
	db     *ydb.Driver
	qTx    query.Transaction
	logger logger.Logger
}

func New(
	db *ydb.Driver,
	logger logger.Logger,
) *PasswordRecoveryRequests {
	return &PasswordRecoveryRequests{
		db:     db,
		logger: logger,
	}
}

func (p *PasswordRecoveryRequests) WithTx(tx query.Transaction) *PasswordRecoveryRequests {
	return &PasswordRecoveryRequests{
		db:     p.db,
		qTx:    tx,
		logger: p.logger,
	}
}

func (p *PasswordRecoveryRequests) executeInQueryTransaction(
	ctx context.Context,
	fn func(c context.Context, tx query.Transaction) error,
) error {
	if p.qTx != nil {
		return fn(ctx, p.qTx)
	}

	return p.db.
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
