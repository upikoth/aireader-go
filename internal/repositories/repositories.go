package repositories

import (
	"context"

	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories/http"
	"github.com/upikoth/aireader-go/internal/repositories/ycp"
	"github.com/upikoth/aireader-go/internal/repositories/ydb"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	YCP  *ycp.Ycp
	YDB  *ydb.YDB
	HTTP *http.HTTP
}

func New(
	log logger.Logger,
	cfg *config.Repositories,
	tp trace.TracerProvider,
) (*Repository, error) {
	ycpInstance, err := ycp.New(log, &cfg.Ycp)

	if err != nil {
		return nil, err
	}

	ydbInstance, err := ydb.New(log, &cfg.Ydb)

	if err != nil {
		return nil, err
	}

	httpInstance, err := http.New(log, &cfg.HTTP, tp)

	if err != nil {
		return nil, err
	}

	return &Repository{
		YCP:  ycpInstance,
		YDB:  ydbInstance,
		HTTP: httpInstance,
	}, nil
}

func (r *Repository) Connect(ctx context.Context) error {
	return r.YDB.Migrate(ctx)
}

func (r *Repository) Disconnect(ctx context.Context) error {
	return r.YDB.Disconnect(ctx)
}
