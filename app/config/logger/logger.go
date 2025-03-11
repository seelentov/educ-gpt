package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
)

var logger *zap.Logger

func Logger() *zap.Logger {
	if logger == nil {
		config := zap.NewProductionConfig()

		if os.Getenv("ENV") == "development" {
			config = zap.NewDevelopmentConfig()
		}

		if os.Getenv("DEBUG") == "true" {
			config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		} else {
			config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		}

		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Encoding = "console"

		filePath := filepath.Join("logs", "app.log")
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}

		fileSyncer := zapcore.AddSync(file)

		multiSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), fileSyncer)

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(config.EncoderConfig),
			multiSyncer,
			config.Level,
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

		logger.Debug("Logger initialized")
	}
	return logger
}
