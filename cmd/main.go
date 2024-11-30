package main

import (
	"ROOmail/config"
	_ "ROOmail/docs"
	"ROOmail/internal/router"
	"ROOmail/pkg/db"
	"ROOmail/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	// Создаем экземпляр логгера
	log := logger.NewLogger()
	log.Info("Приложение запущено", zap.String("version", "1.0.0"))

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация базы данных
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Error("Failed to initialize database", zap.Error(err))
		return
	}

	// Инициализация маршрутизатора
	r := router.NewRouter(database, cfg)

	// Пути к сертификату и ключу
	certFile := `./sertificate/server.crt`
	keyFile := `./sertificate/server.key`

	// Запуск HTTPS-сервера
	log.Info("Server started", zap.String("address", cfg.ServerAddress))
	err = http.ListenAndServeTLS(cfg.ServerAddress, certFile, keyFile, r)
	if err != nil {
		log.Error("Failed to start HTTPS server", zap.Error(err))
	}
}
