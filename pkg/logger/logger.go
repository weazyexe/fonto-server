package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.SugaredLogger

func InitializeLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()
	Zap = logger.Sugar()
	Zap.Info("Logger has been initialized!")
}
