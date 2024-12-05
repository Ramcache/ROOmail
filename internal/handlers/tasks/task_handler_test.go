package tasks_test

import (
	"ROOmail/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ROOmail/internal/handlers/tasks"
	"ROOmail/internal/models"
	"ROOmail/pkg/utils/jwt_token"
	"github.com/stretchr/testify/assert"
)

type TaskService struct{}

func (s *TaskService) CreateTask(ctx context.Context, title, description, dueDateStr, priority string, userIDs []int, filePath string, createdBy int) (string, error) {
	return "1", nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID int, title, description, dueDateStr, priority string, UserIDs []int, currentUserID int) error {
	return nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskID int) (*models.Task, error) {
	return nil, nil
}

func (s *TaskService) GetTasks(ctx context.Context, userID int) ([]models.Task, error) {
	return nil, nil
}

func (s *TaskService) GetTasksByUser(ctx context.Context, userID int) ([]models.Task, error) {
	return nil, nil
}

func (s *TaskService) PatchTask(ctx context.Context, taskID int, updates map[string]interface{}) error {
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID int) error {
	return nil
}

func TestCreateTaskHandler(t *testing.T) {

	logger := logger.NewZapLogger()

	handler := &tasks.TaskHandler{
		Service: &TaskService{},
		Log:     logger,
	}

	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		DueDate:     "2024-12-31",
		Priority:    "High",
		UserIDs:     []int{1, 2, 3},
		FilePath:    "/some/path/file.txt",
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Не удалось сериализовать задачу: %v", err)
	}

	req, err := http.NewRequest("POST", "/admin/tasks/create", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatalf("Не удалось создать HTTP-запрос: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	userClaims := &jwt_token.Claims{
		UserID: 1,
	}
	ctx := context.WithValue(req.Context(), "user", userClaims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.CreateTaskHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	expectedResponse := fmt.Sprintf(`{"message": "Задача успешно создана", "task_id": "1"}`)
	assert.JSONEq(t, expectedResponse, rr.Body.String())
}
