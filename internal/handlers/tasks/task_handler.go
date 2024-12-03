package tasks

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

type TaskHandler struct {
	service *TaskService
	log     logger.Logger
}

func NewTaskHandler(service *TaskService, log logger.Logger) *TaskHandler {
	return &TaskHandler{service: service,
		log: log,
	}
}

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	taskID, err := h.service.CreateTask(r.Context(), req.Title, req.Description, req.DueDate, req.Priority, req.UserIDs, req.FilePath)
	if err != nil {
		h.log.Error("Failed to create task: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Task created successfully", "task_id": %s}`, taskID)))
}
