package users

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/users"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type usersRepositories struct {
	dbDriver *ydb.Driver
	users    *users.Users
}

type Users struct {
	logger       logger.Logger
	repositories *usersRepositories
}

func New(
	logger logger.Logger,
	db *ydb.Driver,
	usersRepo *users.Users,
) *Users {
	return &Users{
		logger: logger,
		repositories: &usersRepositories{
			dbDriver: db,
			users:    usersRepo,
		},
	}
}

func (u *Users) Transaction(
	ctx context.Context,
	fn func(ctxTx context.Context, uTx *Users) error,
	opts ...query.TransactionOption,
) error {
	return u.repositories.dbDriver.Query().Do(ctx, func(ctx context.Context, s query.Session) error {
		tx, err := s.Begin(ctx, query.TxSettings(opts...))

		if err != nil {
			return errors.WithStack(err)
		}

		defer func() { _ = tx.Rollback(ctx) }()

		uTx := u.WithTx(tx)
		if fnErr := fn(ctx, uTx); fnErr != nil {
			return fnErr
		}

		return tx.CommitTx(ctx)
	})
}

func (u *Users) WithTx(tx query.Transaction) *Users {
	return &Users{
		logger:       u.logger,
		repositories: u.repositories.withTx(tx),
	}
}

func (ur *usersRepositories) withTx(tx query.Transaction) *usersRepositories {
	return &usersRepositories{
		dbDriver: ur.dbDriver,
		users:    ur.users.WithTx(tx),
	}
}
