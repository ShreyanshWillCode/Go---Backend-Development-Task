package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global *zap.Logger

func Init(env string) error {
	var (
		cfg	zap.Config
		err	error
	)

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Development = true
	}

	global, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	return nil
}

func Sync() {
	if global != nil {
		_ = global.Sync()
	}
}

func Info(msg string, fields ...zap.Field) {
	global.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	global.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	global.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	global.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	global.Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return global.With(fields...)
}
