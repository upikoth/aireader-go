package ydb

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	passwordrecoveryrequests "github.com/upikoth/aireader-go/internal/repositories/ydb/password-recovery-requests"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/registrations"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/sessions"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/users"
	"github.com/upikoth/aireader-go/internal/repositories/ydb/voices"
	ydbOtel "github.com/ydb-platform/ydb-go-sdk-otel"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/trace"
	yc "github.com/ydb-platform/ydb-go-yc"
)

type YDB struct {
	Users                    *users.Users
	Sessions                 *sessions.Sessions
	Registrations            *registrations.Registrations
	PasswordRecoveryRequests *passwordrecoveryrequests.PasswordRecoveryRequests
	Voices                   *voices.Voices
	DB                       *ydb.Driver
	config                   *config.Ydb
	logger                   logger.Logger
}

func New(
	log logger.Logger,
	cfg *config.Ydb,
) (*YDB, error) {
	filePath := fmt.Sprintf("%s/%s", cfg.AuthFileDirName, cfg.AuthFileName)
	err := writeCredentialsToFile(cfg.AuthFileDirName, cfg.AuthFileName, cfg.YcSaJSONCredentials)

	if err != nil {
		return nil, err
	}

	timeoutSecondsCount := 30
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSecondsCount))
	defer cancel()

	driver, err := ydb.Open(
		ctx,
		cfg.Dsn,
		yc.WithServiceAccountKeyFileCredentials(filePath),
		ydbOtel.WithTraces(
			ydbOtel.WithDetails(trace.DetailsAll),
		),
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &YDB{
		config:                   cfg,
		logger:                   log,
		DB:                       driver,
		Users:                    users.New(driver, log),
		Sessions:                 sessions.New(driver, log),
		Registrations:            registrations.New(driver, log),
		PasswordRecoveryRequests: passwordrecoveryrequests.New(driver, log),
		Voices:                   voices.New(driver, log),
	}, nil
}

//go:embed 1_migrations/*.sql
var embedMigrations embed.FS

func (y *YDB) Migrate(ctx context.Context) error {
	connector, cErr := ydb.Connector(
		y.DB,
		ydb.WithDefaultQueryMode(ydb.ScriptingQueryMode),
		ydb.WithFakeTx(ydb.ScriptingQueryMode),
		ydb.WithAutoDeclare(),
		ydb.WithNumericArgs(),
	)
	if cErr != nil {
		return errors.WithStack(cErr)
	}

	sqlDB := sql.OpenDB(connector)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("ydb"); err != nil {
		return errors.WithStack(err)
	}

	if err := goose.UpContext(ctx, sqlDB, "1_migrations"); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (y *YDB) Disconnect(ctx context.Context) error {
	if y.DB == nil {
		return nil
	}

	return errors.WithStack(y.DB.Close(ctx))
}

func writeCredentialsToFile(dirName string, fileName string, credentials []byte) error {
	filePath := fmt.Sprintf("%s/%s", dirName, fileName)

	if len(credentials) > 0 {
		_, err := os.Stat(dirName)

		if err != nil {
			mkdirErr := os.Mkdir(dirName, 0777)

			if mkdirErr != nil {
				return errors.WithStack(mkdirErr)
			}
		}

		err = os.WriteFile(filePath, credentials, 0600)

		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
