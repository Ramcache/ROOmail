package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

type zapLogger struct {
	sugar *zap.SugaredLogger
}

func NewLogger() Logger {
	logger, _ := zap.NewProduction() // или zap.NewDevelopment(), если требуется больше информации для отладки
	defer logger.Sync()              // гарантируем запись всех логов перед завершением работы программы
	return &zapLogger{
		sugar: logger.Sugar(),
	}
}

func (l *zapLogger) Info(msg string, fields ...interface{}) {
	l.sugar.Infow(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...interface{}) {
	l.sugar.Errorw(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...interface{}) {
	l.sugar.Debugw(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...interface{}) {
	l.sugar.Warnw(msg, fields...)
}
