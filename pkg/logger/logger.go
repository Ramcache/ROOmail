package logger

import (
	"go.uber.org/zap"
)

// Logger - интерфейс для логгера
type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
}

// zapLogger - структура, реализующая интерфейс Logger с использованием zap
type zapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger создает новый экземпляр zapLogger
func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	return &zapLogger{logger: sugar}
}

// Info - обычный лог уровня Info
func (z *zapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

// Infof - форматированный лог уровня Info
func (z *zapLogger) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}

// Error - обычный лог уровня Error
func (z *zapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

// Errorf - форматированный лог уровня Error
func (z *zapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}

// Fatal - обычный лог уровня Fatal (завершает выполнение)
func (z *zapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

// Fatalf - форматированный лог уровня Fatal (завершает выполнение)
func (z *zapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Fatalf(template, args...)
}

// Warn - обычный лог уровня Warn
func (z *zapLogger) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

// Warnf - форматированный лог уровня Warn
func (z *zapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}
