package tasks

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type TaskInterface interface {
	CreateTask(ctx context.Context, title, description, dueDateStr, priority string, userIDs []int, filePath string, createdBy int) (string, error)
}

type TaskService struct {
	db *pgxpool.Pool
}

func NewTaskService(db *pgxpool.Pool) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description, dueDateStr, priority string, userIDs []int, filePath string, createdBy int) (string, error) {
	if title == "" || description == "" {
		return "", fmt.Errorf("Title and description are required")
	}

	var dueDate *time.Time
	if dueDateStr != "" {
		parsedDueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			return "", fmt.Errorf("Invalid due date format")
		}
		dueDate = &parsedDueDate
	}

	var taskID int
	query := `INSERT INTO tasks (title, description, due_date, priority, file_path, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := s.db.QueryRow(ctx, query, title, description, dueDate, priority, filePath, createdBy).Scan(&taskID)
	if err != nil {
		return "", fmt.Errorf("Failed to create task: %w", err)
	}

	for _, userID := range userIDs {
		_, err = s.db.Exec(ctx, `INSERT INTO tasks_users (task_id, user_id, sent_by) VALUES ($1, $2, $3)`, taskID, userID, createdBy)
		if err != nil {
			return "", fmt.Errorf("Failed to assign task to users: %w", err)
		}
	}

	return strconv.Itoa(taskID), nil
}
