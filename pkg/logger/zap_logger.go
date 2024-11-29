package logger

import (
	"go.uber.org/zap"
)

var logInstance *zap.SugaredLogger

// Logger интерфейс для логирования
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// zapLogger реализация интерфейса Logger
type zapLogger struct {
	logger *zap.SugaredLogger
}

// InitLogger инициализирует глобальный логгер
func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	logInstance = logger.Sugar()
}

// GetLogger возвращает глобальный логгер
func GetLogger() *zap.SugaredLogger {
	if logInstance == nil {
		panic("Логгер не инициализирован. Вызовите InitLogger() перед использованием.")
	}
	return logInstance
}

// NewZapLogger создаёт новый экземпляр интерфейса Logger (если нужен отдельный экземпляр)
func NewZapLogger() Logger {
	return &zapLogger{logger: GetLogger()}
}

// Реализация интерфейса Logger
func (z *zapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z *zapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z *zapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}
