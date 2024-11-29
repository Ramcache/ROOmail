package tasks

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var log = logger.GetLogger()

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
	return &TaskHandler{service: service}
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

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		log.Warn("Неавторизованный запрос: ", err)
		utils.RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}

	tasks, err := h.service.GetTaskForUser(userID, schoolID, dueDate)
	if err != nil {
		log.Error("Ошибка получения задач для пользователя: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении списка задач")
		return
	}

	log.Info("Список задач успешно получен для пользователя: ", userID)
	utils.RespondJSON(w, http.StatusOK, tasks)
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
		log.Warn("Некорректный запрос при создании задачи: ", err)
		utils.RespondJSON(w, http.StatusBadRequest, "Некорректный запрос")
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		log.Warn("Неавторизованный запрос при создании задачи: ", err)
		utils.RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	newTask.UserID = userID

	if newTask.Title == "" {
		log.Warn("Пустое название задачи при создании")
		utils.RespondJSON(w, http.StatusBadRequest, "Название обязательно")
		return
	}

	err = h.service.CreateTask(&newTask)
	if err != nil {
		log.Error("Ошибка при сохранении задачи: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при сохранении задачи")
		return
	}

	log.Info("Задача успешно создана для пользователя: ", userID)
	utils.RespondJSON(w, http.StatusCreated, newTask)
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

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			log.Warn("Задача с ID не найдена: ", id)
			utils.RespondJSON(w, http.StatusNotFound, fmt.Sprintf("Задача с ID %s не найдена", id))
			return
		}
		log.Error("Ошибка при получении задачи: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		log.Warn("Неавторизованный запрос на получение задачи: ", err)
		utils.RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	if task.UserID != userID {
		log.Warn("Пользователь с ID ", userID, " пытался получить доступ к задаче, принадлежащей другому пользователю")
		utils.RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	log.Info("Задача с ID успешно получена для пользователя: ", id, userID)
	utils.RespondJSON(w, http.StatusOK, task)
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

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			log.Warn("Задача с ID не найдена для обновления: ", id)
			utils.RespondJSON(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		log.Error("Ошибка при получении задачи для обновления: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		log.Warn("Неавторизованный запрос на обновление задачи: ", err)
		utils.RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	if task.UserID != userID {
		log.Warn("Пользователь с ID ", userID, " пытался обновить задачу, принадлежащую другому пользователю")
		utils.RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	var updatedTask models.Task
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updatedTask)
	if err != nil {
		log.Warn("Некорректный запрос при обновлении задачи: ", err)
		utils.RespondJSON(w, http.StatusBadRequest, "Некорректный запрос")
		return
	}

	err = h.service.UpdateTaskInDB(id, updatedTask)
	if err != nil {
		log.Error("Ошибка при обновлении задачи в базе данных: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при обновлении задачи")
		return
	}

	log.Info("Задача с ID успешно обновлена: ", id)
	utils.RespondJSON(w, http.StatusOK, updatedTask)
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

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			log.Warn("Задача с ID не найдена для удаления: ", id)
			utils.RespondJSON(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		log.Error("Ошибка при получении задачи для удаления: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		log.Warn("Неавторизованный запрос на удаление задачи: ", err)
		utils.RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}
	if task.UserID != userID {
		log.Warn("Пользователь с ID ", userID, " пытался удалить задачу, принадлежащую другому пользователю")
		utils.RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		log.Error("Ошибка при удалении задачи: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при удалении задачи")
		return
	}

	log.Info("Задача с ID успешно удалена: ", id)
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
		log.Warn("Неавторизованный запрос на загрузку файлов: ", err)
		utils.RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}

	formdata := r.MultipartForm
	files := formdata.File["files"]

	err = h.service.UploadFilesForUser(files, userID)
	if err != nil {
		log.Error("Ошибка при загрузке файлов для пользователя: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при обработке файлов: "+err.Error())
		return
	}

	log.Info("Файлы успешно загружены для пользователя: ", userID)
	utils.RespondJSON(w, http.StatusOK, "Файлы успешно загружены")
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
		log.Warn("Неавторизованный запрос на скачивание файла: ", err)
		utils.RespondJSON(w, http.StatusUnauthorized, fmt.Sprintf("Неавторизовано: %v", err))
		return
	}

	task, err := h.service.GetTaskByFileID(fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn("Файл с ID не найден: ", fileID)
			utils.RespondJSON(w, http.StatusNotFound, "Файл не найден")
			return
		}
		log.Error("Ошибка при проверке задачи для файла: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при проверке задачи для файла: "+err.Error())
		return
	}

	if task.UserID != userID {
		log.Warn("Пользователь с ID ", userID, " пытался скачать файл, принадлежащий другому пользователю")
		utils.RespondJSON(w, http.StatusForbidden, "Доступ запрещён")
		return
	}

	// Получение пути к файлу
	filePath, err := h.service.FetchFilePath(fileID)
	if err != nil {
		log.Error("Ошибка при получении пути к файлу: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при получении пути к файлу: "+err.Error())
		return
	}

	err = h.service.ServeFile(w, filePath)
	if err != nil {
		log.Error("Ошибка при скачивании файла: ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при скачивании файла: "+err.Error())
	}
}
