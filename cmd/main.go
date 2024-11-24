package main

import (
	"ROOmail/config"
	"ROOmail/internal/router"
	"ROOmail/pkg/db"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация базы данных
	err = db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Инициализация маршрутизатора
	r := router.NewRouter()

	// Запуск сервера
	log.Printf("Server started at %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
