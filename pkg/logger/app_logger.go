package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var aLogger *zap.Logger

func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	aLogger, _ = config.Build()
}

func Error(message string, err error) {
	aLogger.Error(message, zap.Error(err))
}
