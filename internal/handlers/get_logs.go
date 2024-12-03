package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

var (
	logger, _ = zap.NewProduction() // Создаем логгер
	logDir    = "./logs"
)

// ListLogsHandler
// @Summary Получить список файлов логов
// @Description Возвращает список файлов, находящихся в директории логов.
// @Tags logs
// @Produce json
// @Success 200 {array} string "Список имен файлов логов"
// @Failure 500 {object} map[string]string "Ошибка при чтении директории"
// @Router /admin/logs/list [get]
func ListLogsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(logDir)
	if err != nil {
		logger.Error("Ошибка при чтении директории", zap.Error(err))
		http.Error(w, fmt.Sprintf("Ошибка при чтении директории: %v", err), http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fileNames); err != nil {
		logger.Error("Ошибка при кодировании JSON", zap.Error(err))
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}

// LogsHandler
// @Summary Получить содержимое файла логов
// @Description Возвращает содержимое указанного файла логов из директории.
// @Tags logs
// @Param filename path string true "Имя файла лога"
// @Produce text/plain
// @Success 200 {string} string "Содержимое файла логов"
// @Failure 400 {object} map[string]string "Имя файла не указано или некорректно"
// @Failure 404 {object} map[string]string "Файл не найден"
// @Failure 500 {object} map[string]string "Ошибка при чтении файла логов"
// @Router /admin/logs/{filename} [get]
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	if filename == "" {
		http.Error(w, "Имя файла не указано", http.StatusBadRequest)
		return
	}

	if filepath.Base(filename) != filename {
		http.Error(w, "Некорректное имя файла", http.StatusBadRequest)
		return
	}

	logFilePath := filepath.Join(logDir, filename)

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	logs, err := os.ReadFile(logFilePath)
	if err != nil {
		logger.Error("Ошибка при чтении файла логов", zap.Error(err))
		http.Error(w, fmt.Sprintf("Ошибка при чтении файла логов: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(logs)
}
