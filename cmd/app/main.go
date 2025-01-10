package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/app"
	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	initCtx := context.Background()
	// Чтение .env файла нужно только при локальной разработке.
	// В других случаях значения переменных окружения уже должны быть установлены.
	// Поэтому ошибку загрузки файла обрабатывать не нужно.
	_ = godotenv.Load()
	loggerInstance := logger.New()

	cfg, err := config.New()
	if err != nil {
		loggerInstance.Fatal(err.Error())
		return
	}

	if cfg == nil {
		loggerInstance.Fatal(errors.New("Некорректная инициализация конфига приложения").Error())
		return
	}

	if cfg.Environment == constants.EnvironmentDevelopment {
		loggerInstance.SetPrettyOutputToConsole()
	}

	tp := initTracing(
		&cfg.Controller.HTTP,
		loggerInstance,
	)

	appInstance, err := app.New(cfg, loggerInstance, tp)
	if err != nil {
		loggerInstance.Fatal(err.Error())
	}

	go func() {
		loggerInstance.Info("Запуск приложения")

		if appErr := appInstance.Start(initCtx); !errors.Is(appErr, http.ErrServerClosed) {
			loggerInstance.Fatal(appErr.Error())
		}

		loggerInstance.Info("Приложение перестало принимать новые запросы")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	timeToStopAppInSeconds := 10
	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(),
		time.Duration(timeToStopAppInSeconds)*time.Second,
	)
	defer shutdownRelease()

	if stopErr := appInstance.Stop(shutdownCtx); stopErr != nil {
		loggerInstance.Fatal(fmt.Sprintf("Не удалось корректно остановить сервер, ошибка: %v", stopErr))
	}

	loggerInstance.Info("Приложение остановлено")
}

func initTracing(
	config *config.ControllerHTTP,
	logger logger.Logger,
) trace.TracerProvider {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                config.SentryDsn,
		Environment:        config.Environment,
		DebugWriter:        os.Stdout,
		EnableTracing:      true,
		AttachStacktrace:   true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
		SampleRate:         1.0,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Не удалось инициализировать Sentry: %v\n", errors.WithStack(err)))
	} else {
		logger.Debug("Sentry успешно инициализирована")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)

	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())
	otel.SetTracerProvider(tp)

	return tp
}
