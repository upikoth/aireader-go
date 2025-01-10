package passwordrecoveryrequests

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	passwordrecoveryrequests "github.com/upikoth/aireader-go/internal/repositories/ydb/password-recovery-requests"
	"github.com/upikoth/aireader-go/internal/services/emails"
	"github.com/upikoth/aireader-go/internal/services/sessions"
	"github.com/upikoth/aireader-go/internal/services/users"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type passwordRecoveryRequestsRepositories struct {
	dbDriver                 *ydb.Driver
	passwordRecoveryRequests *passwordrecoveryrequests.PasswordRecoveryRequests
}

type passwordRecoveryRequestsServices struct {
	users    *users.Users
	sessions *sessions.Sessions
	emails   *emails.Emails
}

type PasswordRecoveryRequests struct {
	logger       logger.Logger
	repositories *passwordRecoveryRequestsRepositories
	services     *passwordRecoveryRequestsServices
}

func New(
	logger logger.Logger,
	db *ydb.Driver,
	passwordRecoveryRequestsRepo *passwordrecoveryrequests.PasswordRecoveryRequests,
	usersSrv *users.Users,
	sessionsSrv *sessions.Sessions,
	emailsSrv *emails.Emails,
) *PasswordRecoveryRequests {
	return &PasswordRecoveryRequests{
		logger: logger,
		repositories: &passwordRecoveryRequestsRepositories{
			dbDriver:                 db,
			passwordRecoveryRequests: passwordRecoveryRequestsRepo,
		},
		services: &passwordRecoveryRequestsServices{
			users:    usersSrv,
			sessions: sessionsSrv,
			emails:   emailsSrv,
		},
	}
}

func (p *PasswordRecoveryRequests) Transaction(
	ctx context.Context,
	fn func(ctxTx context.Context, rTx *PasswordRecoveryRequests) error,
	opts ...query.TransactionOption,
) error {
	return p.repositories.dbDriver.Query().Do(ctx, func(ctx context.Context, s query.Session) error {
		tx, err := s.Begin(ctx, query.TxSettings(opts...))

		if err != nil {
			return errors.WithStack(err)
		}

		defer func() { _ = tx.Rollback(ctx) }()

		rTx := p.WithTx(tx)
		if fnErr := fn(ctx, rTx); fnErr != nil {
			return fnErr
		}

		return tx.CommitTx(ctx)
	})
}

func (p *PasswordRecoveryRequests) WithTx(tx query.Transaction) *PasswordRecoveryRequests {
	return &PasswordRecoveryRequests{
		logger:       p.logger,
		repositories: p.repositories.withTx(tx),
		services:     p.services.withTx(tx),
	}
}

func (pr *passwordRecoveryRequestsRepositories) withTx(tx query.Transaction) *passwordRecoveryRequestsRepositories {
	return &passwordRecoveryRequestsRepositories{
		dbDriver:                 pr.dbDriver,
		passwordRecoveryRequests: pr.passwordRecoveryRequests.WithTx(tx),
	}
}

func (ps *passwordRecoveryRequestsServices) withTx(tx query.Transaction) *passwordRecoveryRequestsServices {
	return &passwordRecoveryRequestsServices{
		users:    ps.users.WithTx(tx),
		sessions: ps.sessions.WithTx(tx),
		emails:   ps.emails,
	}
}
