package main

import (
	"ROOmail/config"
	_ "ROOmail/docs"
	"ROOmail/internal/router"
	"ROOmail/pkg/db"
	"ROOmail/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// @title Документация ROOmail API
// @version 1.0
// @description Это документация API для проекта ROOmail.
// @termsOfService http://swagger.io/terms/

// @contact.name Поддержка API
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	// Инициализация логгера
	log := logger.NewZapLogger()
	log.Info("Приложение запущено")

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация базы данных
	dbErr := db.InitDB()
	if dbErr != nil {
		log.Error("Failed to initialize database: ", dbErr)
		os.Exit(1)
	}

	database := db.DB // Используем глобальную переменную DB из пакета db

	// Инициализация маршрутизатора
	r := router.InitRouter(database, cfg)

	// Пути к сертификату и ключу

	// Запуск HTTPS-сервера
	serverAddr := "https://localhost" + cfg.ServerAddress
	log.Infof("Server started at %s", serverAddr)

	err := http.ListenAndServeTLS(cfg.ServerAddress, os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"), r)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}
