package logging

import (
	"tarantool-app/internal/constants"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	SugaredLogger *zap.SugaredLogger
}

// Creates a logger based on application environment.
// Function createLogger does all configuration work.
func NewLogger(env string) *Logger {
	logger := createLogger(env)
	return &Logger{SugaredLogger: logger.Sugar()}
}

// Configures logger based on application environment.
func createLogger(env string) *zap.Logger {
	// Production configuration.
	prodEncoderCfg := zap.NewProductionEncoderConfig()
	prodEncoderCfg.TimeKey = "timestamp"
	prodEncoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	prodConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     prodEncoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	// Local development configuration.
	develEncodefCfg := zap.NewDevelopmentEncoderConfig()
	develEncodefCfg.TimeKey = "timestamp"
	develEncodefCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	develConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     develEncodefCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	var zlogger *zap.Logger

	switch env {
	case constants.EnvLocal:
		zlogger = zap.Must(develConfig.Build())
	case constants.EnvProd:
		zlogger = zap.Must(prodConfig.Build())
	}

	return zlogger
}

// Logging methods wrappers.

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}
