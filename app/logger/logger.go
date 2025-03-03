package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var (
	ErrFailedInit = errors.New("failed to initialize logger")
)

var logger *zap.Logger

func Logger() *zap.Logger {
	if logger == nil {
		config := zap.NewProductionConfig()
		config.Encoding = "json"
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

		l, err := config.Build()
		if err != nil {
			log.Fatalf("%v: %v", ErrFailedInit, err)
		}

		logger = l
		l.Debug("Logger initialized")
	}
	return logger
}
