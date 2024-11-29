package tasks

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// RespondJSON is a helper function to respond with JSON data
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// GetTasksHandler получает список задач
// @Summary Получение списка задач
// @Description Получить список задач по school_id и due_date
// @Tags tasks
// @Accept json
// @Produce json
// @Param school_id query string false "ID школы"
// @Param due_date query string false "Срок выполнения задачи (формат: YYYY-MM-DD)"
// @Success 200 {array} models.Task
// @Failure 500 {object} string "Ошибка при получении списка задач"
// @Router /tasks [get]
func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	schoolID := queryValues.Get("school_id")
	dueDate := queryValues.Get("due_date")

	// Извлечение userID из токена
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		// Использование ошибки для формирования сообщения
		RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}

	// Извлечение задач для конкретного пользователя
	tasks, err := h.service.GetTaskForUser(userID, schoolID, dueDate)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении списка задач")
		return
	}

	// Ответ с задачами
	RespondJSON(w, http.StatusOK, tasks)
}

// CreateTaskHandler создает новую задачу
// @Summary Создание задачи
// @Description Создать новую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Данные новой задачи"
// @Success 201 {object} models.Task
// @Failure 400 {object} string "Некорректный запрос"
// @Failure 500 {object} string "Ошибка при сохранении задачи"
// @Router /tasks [post]
func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newTask)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, "Некорректный запрос")
		return
	}

	// Извлечение user_id из токена
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	newTask.UserID = userID

	if newTask.Title == "" {
		RespondJSON(w, http.StatusBadRequest, "Название обязательно")
		return
	}

	// Создание задачи через сервис
	err = h.service.CreateTask(&newTask)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при сохранении задачи")
		return
	}

	// Возврат созданной задачи с присвоенным идентификатором
	RespondJSON(w, http.StatusCreated, newTask)
}

// GetTaskByIDHandler получает задачу по ID
// @Summary Получение задачи по ID
// @Description Получить задачу по ее уникальному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Success 200 {object} models.Task
// @Failure 404 {object} string "Задача не найдена"
// @Failure 500 {object} string "Ошибка при получении задачи"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Получение задачи из базы
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			RespondJSON(w, http.StatusNotFound, fmt.Sprintf("Задача с ID %s не найдена", id))
			return
		}
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	if task.UserID != userID {
		RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	RespondJSON(w, http.StatusOK, task)
}

// UpdateTaskHandler обновляет задачу
// @Summary Обновление задачи
// @Description Обновить существующую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Param task body models.Task true "Обновленные данные задачи"
// @Success 200 {object} models.Task
// @Failure 400 {object} string "Некорректный запрос"
// @Failure 500 {object} string "Ошибка при обновлении задачи"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Проверка существования задачи
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			RespondJSON(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		return
	}

	// Проверка userID
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	if task.UserID != userID {
		RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	// Декодирование обновленных данных
	var updatedTask models.Task
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updatedTask)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, "Некорректный запрос")
		return
	}

	err = h.service.UpdateTaskInDB(id, updatedTask)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при обновлении задачи")
		return
	}

	RespondJSON(w, http.StatusOK, updatedTask)
}

// DeleteTaskHandler удаляет задачу
// @Summary Удаление задачи
// @Description Удалить задачу по ее уникальному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Success 204 "Задача успешно удалена"
// @Failure 500 {object} string "Ошибка при удалении задачи"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Проверка существования задачи
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			RespondJSON(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		return
	}

	// Проверка userID
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	if task.UserID != userID {
		RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		log.Printf("Ошибка при удалении задачи: %v", err)
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при удалении задачи")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UploadFileHandler загружает файлы
// @Summary Загрузка файлов
// @Description Загрузить файлы
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Файлы для загрузки"
// @Success 200 "Файлы успешно загружены"
// @Failure 500 {object} string "Ошибка при обработке файлов"
// @Router /tasks/upload [post]
func (h *TaskHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}

	formdata := r.MultipartForm
	files := formdata.File["files"]

	err = h.service.UploadFilesForUser(files, userID)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при обработке файлов: "+err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, "Файлы успешно загружены")
}

// DownloadFileHandler скачивает файл
// @Summary Скачивание файла
// @Description Скачать файл по его уникальному идентификатору
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Param fileID path string true "ID файла"
// @Success 200 "Файл успешно скачан"
// @Failure 404 {object} string "Файл не найден"
// @Failure 500 {object} string "Ошибка при скачивании файла"
// @Router /tasks/download/{fileID} [get]
func (h *TaskHandler) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileID"]

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}

	// Проверка, привязан ли файл к задаче пользователя
	task, err := h.service.GetTaskByFileID(fileID) // Новый метод в сервисе
	if err != nil {
		if err == sql.ErrNoRows {
			RespondJSON(w, http.StatusNotFound, "Файл не найден")
			return
		}
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при проверке задачи для файла: "+err.Error())
		return
	}

	if task.UserID != userID {
		RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	// Получение пути к файлу
	filePath, err := h.service.FetchFilePath(fileID)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении пути к файлу: "+err.Error())
		return
	}

	err = h.service.ServeFile(w, filePath)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Ошибка при скачивании файла: "+err.Error())
	}
}
