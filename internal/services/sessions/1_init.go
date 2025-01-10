package sessions

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/sessions"
	"github.com/upikoth/aireader-go/internal/services/users"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type sessionsRepositories struct {
	dbDriver *ydb.Driver
	sessions *sessions.Sessions
}

type sessionsServices struct {
	users *users.Users
}

type Sessions struct {
	logger       logger.Logger
	repositories *sessionsRepositories
	services     *sessionsServices
}

func New(
	log logger.Logger,
	db *ydb.Driver,
	sessionsRepo *sessions.Sessions,
	usersSrv *users.Users,
) *Sessions {
	return &Sessions{
		logger: log,
		repositories: &sessionsRepositories{
			dbDriver: db,
			sessions: sessionsRepo,
		},
		services: &sessionsServices{
			users: usersSrv,
		},
	}
}

func (s *Sessions) Transaction(
	ctx context.Context,
	fn func(ctxTx context.Context, sTx *Sessions) error,
	opts ...query.TransactionOption,
) error {
	return s.repositories.dbDriver.Query().Do(ctx, func(ctx context.Context, qs query.Session) error {
		tx, err := qs.Begin(ctx, query.TxSettings(opts...))

		if err != nil {
			return errors.WithStack(err)
		}

		defer func() { _ = tx.Rollback(ctx) }()

		sTx := s.WithTx(tx)
		if fnErr := fn(ctx, sTx); fnErr != nil {
			return fnErr
		}

		return tx.CommitTx(ctx)
	})
}

func (s *Sessions) WithTx(tx query.Transaction) *Sessions {
	return &Sessions{
		logger:       s.logger,
		repositories: s.repositories.withTx(tx),
		services:     s.services.withTx(tx),
	}
}

func (rr *sessionsRepositories) withTx(tx query.Transaction) *sessionsRepositories {
	return &sessionsRepositories{
		dbDriver: rr.dbDriver,
		sessions: rr.sessions.WithTx(tx),
	}
}

func (rs *sessionsServices) withTx(tx query.Transaction) *sessionsServices {
	return &sessionsServices{
		users: rs.users.WithTx(tx),
	}
}
