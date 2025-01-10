package registrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/registrations"
	"github.com/upikoth/aireader-go/internal/services/emails"
	"github.com/upikoth/aireader-go/internal/services/sessions"
	"github.com/upikoth/aireader-go/internal/services/users"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type registrationsRepositories struct {
	dbDriver      *ydb.Driver
	registrations *registrations.Registrations
}

type registrationsServices struct {
	users    *users.Users
	sessions *sessions.Sessions
	emails   *emails.Emails
}

type Registrations struct {
	logger       logger.Logger
	repositories *registrationsRepositories
	services     *registrationsServices
}

func New(
	log logger.Logger,
	db *ydb.Driver,
	registrationsRepo *registrations.Registrations,
	usersSrv *users.Users,
	sessionsSrv *sessions.Sessions,
	emailsSrv *emails.Emails,
) *Registrations {
	return &Registrations{
		logger: log,
		repositories: &registrationsRepositories{
			dbDriver:      db,
			registrations: registrationsRepo,
		},
		services: &registrationsServices{
			users:    usersSrv,
			sessions: sessionsSrv,
			emails:   emailsSrv,
		},
	}
}

func (r *Registrations) Transaction(
	ctx context.Context,
	fn func(ctxTx context.Context, rTx *Registrations) error,
	opts ...query.TransactionOption,
) error {
	return r.repositories.dbDriver.Query().Do(ctx, func(ctx context.Context, s query.Session) error {
		tx, err := s.Begin(ctx, query.TxSettings(opts...))

		if err != nil {
			return errors.WithStack(err)
		}

		defer func() { _ = tx.Rollback(ctx) }()

		rTx := r.WithTx(tx)
		if fnErr := fn(ctx, rTx); fnErr != nil {
			return fnErr
		}

		return tx.CommitTx(ctx)
	})
}

func (r *Registrations) WithTx(tx query.Transaction) *Registrations {
	return &Registrations{
		logger:       r.logger,
		repositories: r.repositories.withTx(tx),
		services:     r.services.withTx(tx),
	}
}

func (rr *registrationsRepositories) withTx(tx query.Transaction) *registrationsRepositories {
	return &registrationsRepositories{
		dbDriver:      rr.dbDriver,
		registrations: rr.registrations.WithTx(tx),
	}
}

func (rs *registrationsServices) withTx(tx query.Transaction) *registrationsServices {
	return &registrationsServices{
		users:    rs.users.WithTx(tx),
		sessions: rs.sessions.WithTx(tx),
		emails:   rs.emails,
	}
}
