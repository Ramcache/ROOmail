package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
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
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
}

// zapLogger - структура, реализующая интерфейс Logger с использованием zap
type zapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger создает новый экземпляр zapLogger
func NewZapLogger() Logger {
	// Получаем write syncer для файла логов с датой
	fileSyncer := getLogWriterWithDate()

	// Создаем write syncer для консоли
	consoleSyncer := zapcore.AddSync(os.Stdout)

	// Создаем энкодер для логов
	encoder := getEncoder()

	// Объединяем консольный и файловый write syncer'ы
	combinedSyncer := zapcore.NewMultiWriteSyncer(consoleSyncer, fileSyncer)

	// Создаем ядро с объединенным syncer'ом
	core := zapcore.NewCore(encoder, combinedSyncer, zapcore.DebugLevel)

	// Создаем логгер
	logger := zap.New(core, zap.AddCaller())
	sugar := logger.Sugar()

	return &zapLogger{logger: sugar}
}

// getEncoder задает формат записи логов
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "msg"
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriterWithDate создает запись в файл для логов, добавляя сегодняшнюю дату в имя файла
func getLogWriterWithDate() zapcore.WriteSyncer {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Получаем сегодняшнюю дату в формате "YYYY-MM-DD"
	currentDate := time.Now().Format("2006-01-02")

	// Создаем файл логов с датой
	logFileName := filepath.Join(logDir, "logs-"+currentDate+".txt")

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(file)
}

// Логирование
func (z *zapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z *zapLogger) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}

func (z *zapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z *zapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}

func (z *zapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

func (z *zapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Fatalf(template, args...)
}

func (z *zapLogger) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

func (z *zapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}

func (z *zapLogger) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

func (z *zapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Debugf(template, args...)
}
