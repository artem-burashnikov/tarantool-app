package logging

import (
	"log/slog"
	"os"
	"tarantool-app/internal/constants"
)

// Basic logging interface.
type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// stdlib structured logger.
type Slogger struct {
	logger *slog.Logger
}

// Creates a logger based on the application environment.
// Chooses a text handler and sets a logging level.
func NewSlogger(env string) *Slogger {
	var log *slog.Logger

	switch env {
	case constants.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case constants.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return &Slogger{logger: log}
}

// Interface implementation

func (l *Slogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Slogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Slogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Slogger) Error(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}
