package oauthyandex

import (
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/config"
	oauthyandex "github.com/upikoth/aireader-go/internal/generated/oauthyandex"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

type OauthYandex struct {
	logger logger.Logger
	client *oauthyandex.Client
}

func New(
	log logger.Logger,
	cfg *config.OauthYandex,
	tp trace.TracerProvider,
) (*OauthYandex, error) {
	client, err := oauthyandex.NewClient(
		cfg.APIURL,
		oauthyandex.WithTracerProvider(tp),
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &OauthYandex{
		logger: log,
		client: client,
	}, nil
}
