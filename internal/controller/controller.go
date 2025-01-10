package controller

import (
	"context"

	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/controller/http"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/services"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	http   *http.HTTP
	logger logger.Logger
	config *config.Controller
}

func New(
	cfg *config.Config,
	log logger.Logger,
	srvs *services.Services,
	tp trace.TracerProvider,
) (*Controller, error) {
	httpInstance, err := http.New(&cfg.Controller.HTTP, log, srvs, tp)

	if err != nil {
		return nil, err
	}

	return &Controller{
		http:   httpInstance,
		logger: log,
		config: &cfg.Controller,
	}, nil
}

func (c *Controller) Start() error {
	return c.http.Start()
}

func (c *Controller) Stop(ctx context.Context) error {
	return c.http.Stop(ctx)
}
