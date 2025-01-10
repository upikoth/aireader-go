package logger

import (
	loggerzerolog "github.com/upikoth/aireader-go/internal/pkg/logger/logger-zerolog"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	SetPrettyOutputToConsole()
}

func New() Logger {
	return loggerzerolog.New()
}
