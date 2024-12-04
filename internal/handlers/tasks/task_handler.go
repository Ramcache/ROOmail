package tasks

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils/jwt_token"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
// @Tags Задачи
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

	userClaims, ok := r.Context().Value("user").(*jwt_token.Claims)
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

// GetUserTasksHandler обрабатывает запрос на получение задач, назначенных авторизованному пользователю.
// @Summary Получить задачи пользователя
// @Description Возвращает список задач, назначенных авторизованному пользователю.
// @Tags Задачи
// @Produce json
// @Success 200 {array} models.Task "Список задач пользователя"
// @Failure 401 {object} string "Неавторизованный доступ"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /user/tasks/all/get [get]
func (h *TaskHandler) GetUserTasksHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на получение задач для авторизованного пользователя")

	userClaims, ok := r.Context().Value("user").(*jwt_token.Claims)
	if !ok || userClaims.UserID == 0 {
		h.log.Error("Попытка неавторизованного доступа")
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	userID := userClaims.UserID
	h.log.Info("Получение задач для пользователя", "userID:", userID)

	tasks, err := h.service.GetTasksByUser(r.Context(), userID)
	if err != nil {
		h.log.Error("Не удалось получить задачи для пользователя", err)
		http.Error(w, "Не удалось получить задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.log.Error("Ошибка кодирования задач в JSON", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

// GetTasksHandler обрабатывает запрос на получение всех задач, назначенных конкретному пользователю.
// @Summary Получить все задачи пользователя
// @Description Получает список задач, назначенных аутентифицированному пользователю.
// @Tags Задачи
// @Accept  json
// @Produce  json
// @Param   Authorization header string true "Bearer токен"
// @Success 200 {array} models.Task "Список задач, назначенных пользователю"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неавторизованный доступ"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /tasks/get/{id} [get]
func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на получение задач пользователя")

	userClaims, ok := r.Context().Value("user").(*jwt_token.Claims)
	if !ok {
		h.log.Error("Попытка неавторизованного доступа")
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	userID := userClaims.UserID
	h.log.Info("Получение задач для пользователя", "userID: ", userID)

	tasks, err := h.service.GetTasks(r.Context(), userID)
	if err != nil {
		h.log.Error("Не удалось получить задачи", err)
		http.Error(w, "Не удалось получить задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.log.Error("Ошибка кодирования задач в JSON", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

// UpdateTaskHandler обновляет информацию о существующей задаче.
// @Summary Обновить задачу
// @Description Обновление информации о задаче, такой как название, описание, срок выполнения, приоритет и список пользователей.
// @Tags Задачи
// @Accept  json
// @Produce  json
// @Param   id    path      int   true  "Идентификатор задачи"
// @Param   task  body      models.Task  true  "Данные задачи для обновления"
// @Success 200 {string} string "{"message": "Задача успешно обновлена"}"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /admin/tasks/update/{id} [put]
func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на обновление задачи")

	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		h.log.Error("Некорректный идентификатор задачи", err)
		http.Error(w, "Некорректный идентификатор задачи", http.StatusBadRequest)
		return
	}

	var req models.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Предоставлен некорректный JSON", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	userClaims, ok := r.Context().Value("user").(*jwt_token.Claims)
	if !ok {
		h.log.Error("Попытка неавторизованного доступа")
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	h.log.Info("Обновление задачи", " обновляется пользователем: ", userClaims.UserID)

	currentUserID := userClaims.UserID
	err = h.service.UpdateTask(r.Context(), taskID, req.Title, req.Description, req.DueDate, req.Priority, req.UserIDs, currentUserID)
	if err != nil {
		h.log.Error("Не удалось обновить задачу", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("Задача успешно обновлена", " taskID: ", taskID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Задача успешно обновлена"}`))
}

// PatchTaskHandler обновляет отдельные поля задачи по её идентификатору
// @Summary Частичное обновление задачи
// @Description Обновление одного или нескольких полей задачи по её идентификатору
// @Tags Задачи
// @Accept  json
// @Produce  json
// @Param id path int true "Идентификатор задачи"
// @Param updates body object true "Обновляемые поля задачи"
// @Success 200 {object} map[string]string "Задача успешно обновлена"
// @Failure 400 {string} string "Некорректный идентификатор задачи или JSON"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /admin/tasks/update/{id} [patch]
func (h *TaskHandler) PatchTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на частичное обновление задачи")

	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		h.log.Error("Некорректный идентификатор задачи", err)
		http.Error(w, "Некорректный идентификатор задачи", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.log.Error("Предоставлен некорректный JSON", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	userClaims, ok := r.Context().Value("user").(*jwt_token.Claims)
	if !ok {
		h.log.Error("Попытка неавторизованного доступа")
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	h.log.Info("Частичное обновление задачи", " обновляется пользователем: ", userClaims.UserID)

	err = h.service.PatchTask(r.Context(), taskID, updates)
	if err != nil {
		h.log.Error("Не удалось обновить задачу", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("Задача успешно обновлена", " taskID: ", taskID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Задача успешно обновлена"}`))
}

// DeleteTaskHandler godoc
// @Summary Удаление задачи
// @Description Удаляет задачу и все связанные с ней данные
// @Tags Задачи
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {string} string "Задача успешно удалена"
// @Failure 400 {string} string "Некорректный идентификатор задачи"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Failure 500 {string} string "Не удалось удалить задачу"
// @Router /admin/tasks/delete/{id} [delete]
func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на удаление задачи")

	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		h.log.Error("Некорректный идентификатор задачи", err)
		http.Error(w, "Некорректный идентификатор задачи", http.StatusBadRequest)
		return
	}

	userClaims, ok := r.Context().Value("user").(*jwt_token.Claims)
	if !ok {
		h.log.Error("Попытка неавторизованного доступа")
		http.Error(w, "Неавторизованный доступ", http.StatusUnauthorized)
		return
	}

	h.log.Info("Удаление задачи", " taskID: ", taskID, " выполняется пользователем: ", userClaims.UserID)

	err = h.service.DeleteTask(r.Context(), taskID)
	if err != nil {
		h.log.Error("Не удалось удалить задачу", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("Задача успешно удалена", " taskID: ", taskID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Задача успешно удалена"}`))
}
