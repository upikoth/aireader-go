package oauthmailru

import (
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/config"
	oauthmailru "github.com/upikoth/aireader-go/internal/generated/oauthmailru"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

type OauthMailRu struct {
	logger logger.Logger
	client *oauthmailru.Client
}

func New(
	log logger.Logger,
	cfg *config.OauthMailRu,
	tp trace.TracerProvider,
) (*OauthMailRu, error) {
	client, err := oauthmailru.NewClient(
		cfg.APIURL,
		oauthmailru.WithTracerProvider(tp),
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &OauthMailRu{
		logger: log,
		client: client,
	}, nil
}
