package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

var DB *sql.DB

func InitDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	DB = db
	return db, nil
}

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, password_hash, role FROM users WHERE username=$1"
	err := DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}
