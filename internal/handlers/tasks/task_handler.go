package tasks

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils/JWT"
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

// CreateTaskHandler создает новую задачу
// @Summary Создание новой задачи
// @Description Создает новую задачу с указанными данными
// @Tags задачи
// @Accept json
// @Produce json
// @Param task body models.Task true "Данные задачи"
// @Success 201 {object} map[string]interface{} "Задача успешно создана"
// @Failure 400 {string} string "Неверный JSON"
// @Failure 401 {string} string "Неавторизован"
// @Failure 500 {string} string "Ошибка создания задачи"
// @Router /admin/tasks/create [post]
func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на создание новой задачи")

	var req models.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Предоставлен некорректный JSON", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	userClaims, ok := r.Context().Value("user").(*JWT.Claims)
	if !ok {
		h.log.Error("Попытка неавторизованного доступа")
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	createdBy := userClaims.UserID
	h.log.Info("Создание задачи", " создано пользователем: ", createdBy)

	taskID, err := h.service.CreateTask(r.Context(), req.Title, req.Description, req.DueDate, req.Priority, req.UserIDs, req.FilePath, createdBy)
	if err != nil {
		h.log.Error("Не удалось создать задачу", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("Задача успешно создана", " taskID: ", taskID)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"message": "Задача успешно создана", "task_id": "%s"}`, taskID)))
}
