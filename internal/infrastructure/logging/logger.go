package logging

import (
	"log/slog"
	"os"
	"tarantool-app/internal/constants"
)

// Basic logging interface.
type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
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
	l.logger.Error(msg, args...)
}

// error attribute for error log level.
func (l *Slogger) Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// key=value attribute for any log level.
func (l *Slogger) Str(key, msg string) slog.Attr {
	return slog.String(key, msg)
}
