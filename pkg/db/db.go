package db

import (
	"ROOmail/internal/models"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
	"log"
	"os"
)

var DB *pgxpool.Pool

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set in environment variables")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	log.Println("Database connection established")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, password_hash, role FROM users WHERE username=$1"
	err := DB.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}
