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
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := config.LoadConfig()

	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := router.NewRouter(database, cfg)

	log.Printf("Server started at %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
