package voices

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/voices"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type voicesRepositories struct {
	dbDriver *ydb.Driver
	voices   *voices.Voices
}

type Voices struct {
	logger       logger.Logger
	repositories *voicesRepositories
}

func New(
	logger logger.Logger,
	db *ydb.Driver,
	voicesRepo *voices.Voices,
) *Voices {
	return &Voices{
		logger: logger,
		repositories: &voicesRepositories{
			dbDriver: db,
			voices:   voicesRepo,
		},
	}
}

func (v *Voices) Transaction(
	ctx context.Context,
	fn func(ctxTx context.Context, vTx *Voices) error,
	opts ...query.TransactionOption,
) error {
	return v.repositories.dbDriver.Query().Do(ctx, func(ctx context.Context, s query.Session) error {
		tx, err := s.Begin(ctx, query.TxSettings(opts...))

		if err != nil {
			return errors.WithStack(err)
		}

		defer func() { _ = tx.Rollback(ctx) }()

		uTx := v.WithTx(tx)
		if fnErr := fn(ctx, uTx); fnErr != nil {
			return fnErr
		}

		return tx.CommitTx(ctx)
	})
}

func (v *Voices) WithTx(tx query.Transaction) *Voices {
	return &Voices{
		logger:       v.logger,
		repositories: v.repositories.withTx(tx),
	}
}

func (vr *voicesRepositories) withTx(tx query.Transaction) *voicesRepositories {
	return &voicesRepositories{
		dbDriver: vr.dbDriver,
		voices:   vr.voices.WithTx(tx),
	}
}
