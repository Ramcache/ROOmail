package main

import (
	"ROOmail/config"
	_ "ROOmail/docs"
	"ROOmail/internal/router"
	"ROOmail/pkg/db"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация базы данных
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Инициализация маршрутизатора
	r := router.NewRouter(database, cfg)

	// Пути к сертификату и ключу
	certFile := `C:\Users\Рамзан\server.crt`
	keyFile := `C:\Users\Рамзан\server.key`

	// Запуск HTTPS-сервера
	log.Printf("Server started at https://%s", cfg.ServerAddress)
	err = http.ListenAndServeTLS(cfg.ServerAddress, certFile, keyFile, r)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}
