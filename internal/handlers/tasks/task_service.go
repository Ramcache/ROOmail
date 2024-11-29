package tasks

import (
	"ROOmail/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

// GetTaskByID retrieves a task by its ID
func (s *TaskService) GetTaskByID(taskID string) (*models.Task, error) {
	var task models.Task
	err := s.db.QueryRow(`
		SELECT id, user_id, title, description, due_date, file, priority, schools
		FROM tasks WHERE id = $1`, taskID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.File,
		&task.Priority,
		&task.Schools,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

// CreateTask inserts a new task into the database
func (s *TaskService) CreateTask(task *models.Task) error {
	query := `INSERT INTO tasks (user_id, title, description, due_date, file, priority, schools) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := s.db.QueryRow(query, task.UserID, task.Title, task.Description, task.DueDate, task.File, task.Priority, task.Schools).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}

// GetTask retrieves tasks by school ID and/or due date
func (s *TaskService) GetTask(schoolID, dueDate string) ([]models.Task, error) {
	query := `SELECT id, user_id, title, description, due_date, file, priority, schools FROM tasks`
	var filters []string
	var args []interface{}

	argIndex := 1
	if schoolID != "" {
		filters = append(filters, fmt.Sprintf(`schools LIKE '%%' || $%d || '%%'`, argIndex))
		args = append(args, schoolID)
		argIndex++
	}
	if dueDate != "" {
		filters = append(filters, fmt.Sprintf(`due_date = $%d`, argIndex))
		args = append(args, dueDate)
		argIndex++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var schools string
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.DueDate, &task.File, &task.Priority, &schools)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		task.Schools = strings.Join(parseSchools(schools), ", ")
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return tasks, nil
}

// GetTaskForUser retrieves tasks for a specific user, optionally filtering by school ID and/or due date
func (s *TaskService) GetTaskForUser(userID int, schoolID, dueDate string) ([]models.Task, error) {
	// Преобразуем userID из int в строку
	userIDStr := strconv.Itoa(userID)

	query := `SELECT id, user_id, title, description, due_date, file, priority, schools FROM tasks WHERE user_id = $1`
	var filters []string
	var args []interface{}

	args = append(args, userIDStr) // userID теперь преобразован в строку

	if schoolID != "" {
		filters = append(filters, `schools LIKE '%' || $2 || '%'`)
		args = append(args, schoolID)
	}
	if dueDate != "" {
		filters = append(filters, `due_date = $3`)
		args = append(args, dueDate)
	}

	if len(filters) > 0 {
		query += " AND " + strings.Join(filters, " AND ")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var schools string
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.DueDate, &task.File, &task.Priority, &schools)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		task.Schools = strings.Join(parseSchools(schools), ", ")
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return tasks, nil
}

// UpdateTaskInDB updates an existing task in the database
func (s *TaskService) UpdateTaskInDB(id string, updatedTask models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, due_date = $3, priority = $4, schools = $5, file = $6 WHERE id = $7`
	_, err := s.db.Exec(query, updatedTask.Title, updatedTask.Description, updatedTask.DueDate, updatedTask.Priority, updatedTask.Schools, updatedTask.File, id)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

// DeleteTask deletes a task from the database
func (s *TaskService) DeleteTask(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// UploadFiles uploads multiple files
func (s *TaskService) UploadFiles(files []*multipart.FileHeader) error {
	for _, fileHeader := range files {
		if err := s.uploadFile(fileHeader); err != nil {
			return err
		}
	}
	return nil
}

// UploadFilesForUser uploads files for a specific user
func (s *TaskService) UploadFilesForUser(files []*multipart.FileHeader, userID int) error {
	uploadDir := fmt.Sprintf("./uploads/user_%d", userID)

	// Создание директории при её отсутствии
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	for _, fileHeader := range files {
		// Добавление уникальности имени файла
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
		filePath := filepath.Join(uploadDir, uniqueName)
		if err := s.uploadFileToPath(fileHeader, filePath); err != nil {
			return err
		}
	}

	return nil
}

// GetTaskByFileID retrieves a task by its associated file ID
func (s *TaskService) GetTaskByFileID(fileID string) (*models.Task, error) {
	query := `
		SELECT t.id, t.title, t.user_id
		FROM tasks t
		JOIN files f ON f.task_id = t.id
		WHERE f.id = $1
	`
	var task models.Task
	err := s.db.QueryRow(query, fileID).Scan(&task.ID, &task.Title, &task.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task by file ID: %w", err)
	}
	return &task, nil
}

// FetchFilePath retrieves the file path by file ID
func (s *TaskService) FetchFilePath(fileID string) (string, error) {
	var filePath string
	err := s.db.QueryRow("SELECT file_path FROM files WHERE id = $1", fileID).Scan(&filePath)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file path: %w", err)
	}
	return filePath, nil
}

// ServeFile serves the file to the HTTP response
func (s *TaskService) ServeFile(w http.ResponseWriter, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Определение MIME-типа
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file for content type detection: %w", err)
	}
	file.Seek(0, io.SeekStart)
	contentType := http.DetectContentType(buffer)

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// Helper function to parse schools
func parseSchools(schools string) []string {
	if schools == "" {
		return []string{}
	}
	return strings.Split(schools, ",")
}

// Helper function to upload a single file
func (s *TaskService) uploadFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("unable to open the file: %w", err)
	}
	defer file.Close()

	filePath := filepath.Join("uploads", fileHeader.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create the file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("unable to save the file: %w", err)
	}

	_, err = s.db.Exec("INSERT INTO files (file_path) VALUES ($1)", filePath)
	if err != nil {
		return fmt.Errorf("unable to save file path to database: %w", err)
	}

	return nil
}

// Helper function to upload a file to a specific path
func (s *TaskService) uploadFileToPath(fileHeader *multipart.FileHeader, filePath string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("unable to open the file: %w", err)
	}
	defer file.Close()

	destFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}

	return nil
}

// splitCommaSeparated converts a comma-separated string to a slice of strings
func splitCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, ",")
}
