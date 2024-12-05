package tasks

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/utils/jwt_token"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, title, description, dueDateStr, priority string, userIDs []int, filePath string, createdBy int) (string, error)
	UpdateTask(ctx context.Context, taskID int, title, description, dueDateStr, priority string, UserIDs []int, currentUserID int) error
	GetTaskByID(ctx context.Context, taskID int) (*models.Task, error)
	GetTasks(ctx context.Context, userID int) ([]models.Task, error)
	GetTasksByUser(ctx context.Context, userID int) ([]models.Task, error)
	PatchTask(ctx context.Context, taskID int, updates map[string]interface{}) error
	DeleteTask(ctx context.Context, taskID int) error
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

func (s *TaskService) GetTaskByID(ctx context.Context, taskID int) (*models.Task, error) {
	query := `
		SELECT id, title, description, due_date, priority, file_path, created_by
		FROM tasks
		WHERE id = $1
	`

	var task models.Task
	var dueDate sql.NullTime

	err := s.db.QueryRow(ctx, query, taskID).Scan(&task.ID, &task.Title, &task.Description, &dueDate, &task.Priority, &task.FilePath, &task.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Task not found")
		}
		return nil, fmt.Errorf("Failed to retrieve task: %w", err)
	}

	if dueDate.Valid {
		task.DueDate = dueDate.Time.Format("2006-01-02")
	} else {
		task.DueDate = ""
	}

	userQuery := `SELECT user_id FROM tasks_users WHERE task_id = $1`
	rows, err := s.db.Query(ctx, userQuery, task.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve users for task %d: %w", task.ID, err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("Failed to scan user_id for task %d: %w", task.ID, err)
		}
		userIDs = append(userIDs, userID)
	}
	task.UserIDs = userIDs

	return &task, nil
}

func (s *TaskService) GetTasksByUser(ctx context.Context, userID int) ([]models.Task, error) {
	query := `
		SELECT t.id, t.title, t.description, t.due_date, t.priority, t.file_path, t.created_by
		FROM tasks t
		JOIN tasks_users tu ON t.id = tu.task_id
		WHERE tu.user_id = $1
		ORDER BY t.due_date ASC
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Не удалось получить задачи для пользователя %d: %w", userID, err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate sql.NullTime

		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &dueDate, &task.Priority, &task.FilePath, &task.CreatedBy); err != nil {
			return nil, fmt.Errorf("Не удалось отсканировать данные задачи: %w", err)
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time.Format("2006-01-02")
		} else {
			task.DueDate = ""
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *TaskService) GetTasks(ctx context.Context, userID int) ([]models.Task, error) {
	fmt.Printf("Retrieving tasks for userID: %d\n", userID)

	query := `
		SELECT t.id, t.title, t.description, t.due_date, t.priority, t.file_path, t.created_by
		FROM tasks t
		JOIN tasks_users tu ON t.id = tu.task_id
		WHERE tu.user_id = $1
		ORDER BY t.due_date ASC
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve tasks for user %d: %w", userID, err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate sql.NullTime

		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &dueDate, &task.Priority, &task.FilePath, &task.CreatedBy); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %w", err)
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time.Format("2006-01-02")
		} else {
			task.DueDate = ""
		}

		var userIDs []int
		userQuery := `SELECT user_id FROM tasks_users WHERE task_id = $1`
		userRows, err := s.db.Query(ctx, userQuery, task.ID)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve users for task %d: %w", task.ID, err)
		}
		defer userRows.Close()

		for userRows.Next() {
			var userID int
			if err := userRows.Scan(&userID); err != nil {
				return nil, fmt.Errorf("Failed to scan user_id for task %d: %w", task.ID, err)
			}
			userIDs = append(userIDs, userID)
		}
		task.UserIDs = userIDs

		tasks = append(tasks, task)
	}

	fmt.Printf("Found %d tasks for userID %d\n", len(tasks), userID)

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID int, title, description, dueDateStr, priority string, UserIDs []int, currentUserID int) error {
	if title == "" || description == "" {
		return fmt.Errorf("Title and description are required")
	}

	var dueDate *time.Time
	if dueDateStr != "" {
		parsedDueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			return fmt.Errorf("Invalid due date format")
		}
		dueDate = &parsedDueDate
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := `UPDATE tasks SET title = $1, description = $2, due_date = $3, priority = $4 WHERE id = $5`
	_, err = tx.Exec(ctx, query, title, description, dueDate, priority, taskID)
	if err != nil {
		fmt.Printf("Error while updating task %d: %v\n", taskID, err)
		return fmt.Errorf("Failed to update task: %w", err)
	}

	var currentUserIDs []int
	rows, err := tx.Query(ctx, `SELECT user_id FROM tasks_users WHERE task_id = $1`, taskID)
	if err != nil {
		fmt.Printf("Error while retrieving current users for task %d: %v\n", taskID, err)
		return fmt.Errorf("Failed to retrieve current users for task: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			fmt.Printf("Error while scanning user_id: %v\n", err)
			return fmt.Errorf("Failed to scan user_id: %w", err)
		}
		currentUserIDs = append(currentUserIDs, userID)
	}

	fmt.Printf("Current users found: %v\n", currentUserIDs)

	toRemove := difference(currentUserIDs, UserIDs)
	toAdd := difference(UserIDs, currentUserIDs)

	fmt.Printf("Users to remove: %v\n", toRemove)
	fmt.Printf("Users to add: %v\n", toAdd)

	for _, userID := range toRemove {
		fmt.Printf("Removing user %d from task %d\n", userID, taskID) // Логируем удаление пользователя
		_, err = tx.Exec(ctx, `DELETE FROM tasks_users WHERE task_id = $1 AND user_id = $2`, taskID, userID)
		if err != nil {
			fmt.Printf("Error while removing user %d: %v\n", userID, err)
			return fmt.Errorf("Failed to unassign task from user %d: %w", userID, err)
		}
	}

	for _, userID := range toAdd {
		fmt.Printf("Adding new user %d to task %d\n", userID, taskID) // Логируем добавление пользователя
		_, err = tx.Exec(ctx, `INSERT INTO tasks_users (task_id, user_id, assigned_at, sent_by) VALUES ($1, $2, NOW(), $3)`, taskID, userID, currentUserID)
		if err != nil {
			fmt.Printf("Error while adding user %d: %v\n", userID, err)
			return fmt.Errorf("Failed to assign task to user %d: %w", userID, err)
		}
	}

	return nil
}

func difference(a, b []int) []int {
	m := make(map[int]struct{}, len(b))
	for _, item := range b {
		m[item] = struct{}{}
	}
	var diff []int
	for _, item := range a {
		if _, found := m[item]; !found {
			diff = append(diff, item)
		}
	}
	return diff
}

func (s *TaskService) PatchTask(ctx context.Context, taskID int, updates map[string]interface{}) error {
	query := "UPDATE tasks SET "

	var params []interface{}
	paramCounter := 1
	hasFieldsToUpdate := false

	for key, value := range updates {
		if key == "user_ids" {
			var userIDs []int
			switch v := value.(type) {
			case []int:
				userIDs = v
			case []interface{}:
				for _, id := range v {
					userID, ok := id.(float64)
					if !ok {
						return fmt.Errorf("Invalid user_ids format")
					}
					userIDs = append(userIDs, int(userID))
				}
			default:
				return fmt.Errorf("Invalid user_ids format")
			}

			rows, err := s.db.Query(ctx, `SELECT user_id, sent_by FROM tasks_users WHERE task_id = $1`, taskID)
			if err != nil {
				return fmt.Errorf("Failed to retrieve current users for task: %w", err)
			}
			defer rows.Close()

			currentSentBy := make(map[int]int)
			for rows.Next() {
				var userID, sentBy int
				if err := rows.Scan(&userID, &sentBy); err != nil {
					return fmt.Errorf("Failed to scan user_id and sent_by: %w", err)
				}
				currentSentBy[userID] = sentBy
			}

			_, err = s.db.Exec(ctx, `DELETE FROM tasks_users WHERE task_id = $1`, taskID)
			if err != nil {
				return fmt.Errorf("Failed to remove current users for task: %w", err)
			}

			for _, userID := range userIDs {
				sentBy, exists := currentSentBy[userID]
				if !exists {
					userClaims, ok := ctx.Value("user").(*jwt_token.Claims)
					if !ok {
						return fmt.Errorf("Failed to retrieve user claims from context")
					}
					sentBy = userClaims.UserID
				}
				_, err = s.db.Exec(ctx, `INSERT INTO tasks_users (task_id, user_id, assigned_at, sent_by) VALUES ($1, $2, NOW(), $3)`, taskID, userID, sentBy)
				if err != nil {
					return fmt.Errorf("Failed to assign task to user %d: %w", userID, err)
				}
			}
			continue
		}

		if paramCounter > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", key, paramCounter)
		params = append(params, value)
		paramCounter++
		hasFieldsToUpdate = true
	}

	if hasFieldsToUpdate {
		query += fmt.Sprintf(" WHERE id = $%d", paramCounter)
		params = append(params, taskID)

		_, err := s.db.Exec(ctx, query, params...)
		if err != nil {
			return fmt.Errorf("Failed to patch task: %w", err)
		}
	} else {
		fmt.Println("No fields to update in tasks, only user_ids updated")
	}

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `DELETE FROM tasks_users WHERE task_id = $1`, taskID)
	if err != nil {
		return fmt.Errorf("Failed to delete task-user associations: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, taskID)
	if err != nil {
		return fmt.Errorf("Failed to delete task: %w", err)
	}

	return nil
}
