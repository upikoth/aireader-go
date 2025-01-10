package app

import (
	"context"
	"fmt"

	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/controller"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories"
	"github.com/upikoth/aireader-go/internal/services"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	config       *config.Config
	logger       logger.Logger
	repositories *repositories.Repository
	services     *services.Services
	controller   *controller.Controller
}

func New(
	cfg *config.Config,
	log logger.Logger,
	tp trace.TracerProvider,
) (*App, error) {
	repositoriesInstance, err := repositories.New(log, &cfg.Repositories, tp)

	if err != nil {
		log.Error(fmt.Sprintf("Ошибка при инициализации repositories: %s", err))
		return nil, err
	}

	servicesInstance, err := services.New(log, cfg, repositoriesInstance)

	if err != nil {
		log.Error(fmt.Sprintf("Ошибка при инициализации services: %s", err))
		return nil, err
	}

	controllerInstance, err := controller.New(cfg, log, servicesInstance, tp)

	if err != nil {
		log.Error(fmt.Sprintf("Ошибка при инициализации controller: %s", err))
		return nil, err
	}

	return &App{
		config:       cfg,
		logger:       log,
		repositories: repositoriesInstance,
		services:     servicesInstance,
		controller:   controllerInstance,
	}, nil
}

func (s *App) Start(ctx context.Context) error {
	err := s.repositories.Connect(ctx)

	if err != nil {
		return err
	}

	s.logger.Debug("Подключение к repositories прошло без ошибок")

	return s.controller.Start()
}

func (s *App) Stop(ctx context.Context) error {
	err := s.repositories.Disconnect(ctx)

	if err != nil {
		return err
	}

	s.logger.Debug("Отключение от repositories прошло без ошибок")

	return s.controller.Stop(ctx)
}
