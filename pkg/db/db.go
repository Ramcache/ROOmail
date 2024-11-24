package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("pgx", dataSourceName)
	if err != nil {
		return err
	}
	return DB.Ping()
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
