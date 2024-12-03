package users

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/utils"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUsersService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) AddUser(ctx context.Context, username, password, role string) (int, error) {
	// Хешируем пароль перед сохранением в базу
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("Failed to hash password: %w", err)
	}

	// SQL-запрос для добавления пользователя
	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id`
	var userID int
	err = s.db.QueryRow(ctx, query, username, passwordHash, role).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("Failed to add user to the database: %w", err)
	}

	return userID, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("Не удалось удалить пользователя с id %d: %w", userID, err)
	}
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID int, username, password, role *string) error {
	query := `UPDATE users SET `
	params := []interface{}{}
	paramIndex := 1

	if username != nil {
		query += fmt.Sprintf("username = $%d, ", paramIndex)
		params = append(params, *username)
		paramIndex++
	}

	if password != nil {
		hashedPassword, err := utils.HashPassword(*password)
		if err != nil {
			return fmt.Errorf("Не удалось хешировать пароль: %w", err)
		}
		query += fmt.Sprintf("password_hash = $%d, ", paramIndex)
		params = append(params, hashedPassword)
		paramIndex++
	}

	if role != nil {
		query += fmt.Sprintf("role = $%d, ", paramIndex)
		params = append(params, *role)
		paramIndex++
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", paramIndex)
	params = append(params, userID)

	_, err := s.db.Exec(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("Не удалось обновить пользователя с ID %d: %w", userID, err)
	}

	return nil
}

func (s *UserService) GetUsers(username string) ([]models.UsersList, error) {
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
