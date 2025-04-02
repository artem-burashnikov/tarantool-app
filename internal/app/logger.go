package app

import (
	"tarantool-app/internal/interfaces"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	SugaredLogger *zap.SugaredLogger
}

var _ interfaces.Logger = ZapLogger{} // ZapLogger must satisfy Logger

func NewLogger(env string) ZapLogger {
	logger := configureLogger(env)
	return ZapLogger{SugaredLogger: logger.Sugar()}
}

func configureLogger(env string) *zap.Logger {
	var loggerConfig zap.Config

	switch env {
	case "development":
		loggerConfig = zap.NewDevelopmentConfig()
	case "production":
		loggerConfig = zap.NewProductionConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.DisableCaller = true

	return zap.Must(loggerConfig.Build())
}

func (l ZapLogger) Info(msg string, keysAndValues ...any) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

func (l ZapLogger) Debug(msg string, keysAndValues ...any) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (l ZapLogger) Warn(msg string, keysAndValues ...any) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (l ZapLogger) Error(msg string, keysAndValues ...any) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (l ZapLogger) Fatal(msg string, keysAndValues ...any) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l ZapLogger) Sync() error {
	err := l.SugaredLogger.Sync()
	if err != nil {
		l.Warn("Logger sync error",
			"error", err,
		)
	}
	return err
}
