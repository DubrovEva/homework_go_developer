package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	globalLogger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Logger() *zap.Logger {
	if globalLogger == nil {
		InitLogger()
	}
	return globalLogger
}

func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Logger()
	}

	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return Logger().With(zap.String("trace_id", traceID))
	}

	return Logger()
}

func Debug(msg string, fields ...zap.Field) {
	Logger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger().Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return Logger().With(fields...)
}

func Sync() error {
	return Logger().Sync()
}
