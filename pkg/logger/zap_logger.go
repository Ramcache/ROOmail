package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	return &zapLogger{logger: sugar}
}

func (z *zapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z *zapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z *zapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}
