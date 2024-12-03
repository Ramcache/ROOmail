package users

import (
	"ROOmail/internal/models"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type UsersService struct {
	db *pgxpool.Pool
}

func NewUsersService(db *pgxpool.Pool) *UsersService {
	return &UsersService{db: db}
}

func (s *UsersService) GetUsers(username string) ([]models.UsersList, error) {
	query := "SELECT id, username FROM users"
	var args []interface{}

	if username != "" {
		query += " WHERE username LIKE ?"
		args = append(args, "%"+username+"%")
	}

	rows, err := s.db.Query(context.Background(), query, args...)
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
