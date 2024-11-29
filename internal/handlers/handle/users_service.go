package handle

import (
	"ROOmail/internal/models"
	"database/sql"
	"fmt"
)

type UsersService struct {
	db *sql.DB
}

func NewUsersService(db *sql.DB) *UsersService {
	return &UsersService{db: db}
}

func (s *UsersService) GetUsers(username string) ([]models.UsersList, error) {
	query := "SELECT id, username FROM users"
	var args []interface{}

	if username != "" {
		query += " WHERE username LIKE ?"
		args = append(args, "%"+username+"%")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var users []models.UsersList
	for rows.Next() {
		var user models.UsersList
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %w", err)
	}

	return users, nil
}
