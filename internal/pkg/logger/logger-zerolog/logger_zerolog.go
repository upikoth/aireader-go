package loggerzerolog

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerZeroLog struct{}

func New() *LoggerZeroLog {
	return &LoggerZeroLog{}
}

func (l *LoggerZeroLog) Debug(msg string) {
	log.Debug().Msg(msg)
}

func (l *LoggerZeroLog) Info(msg string) {
	log.Info().Msg(msg)
}

func (l *LoggerZeroLog) Warn(msg string) {
	log.Warn().Msg(msg)
}

func (l *LoggerZeroLog) Error(msg string) {
	log.Error().Msg(msg)
}

func (l *LoggerZeroLog) Fatal(msg string) {
	log.Fatal().Msg(msg)
}

func (l *LoggerZeroLog) SetPrettyOutputToConsole() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
