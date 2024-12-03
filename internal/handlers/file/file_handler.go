package file

import (
	"ROOmail/pkg/logger"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type FileHandler struct {
	service *FileService
	log     logger.Logger
}

func NewFileHandler(service *FileService, log logger.Logger) *FileHandler {
	return &FileHandler{service: service,
		log: log,
	}
}

// UploadFileHandler godoc
// @Summary Загрузка файла
// @Description Загрузка файла на сервер
// @Tags файлы
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]string "{"file_path": "uploaded/file/path"}"
// @Failure 400 {string} string "Ошибка разбора формы"
// @Failure 500 {string} string "Ошибка чтения файла или сохранения файла"
// @Router /users/files/upload [post]
func (h *FileHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на загрузку файла")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.log.Error("Ошибка разбора формы", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.log.Error("Ошибка чтения файла", err)
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filePath, err := h.service.SaveFile(file, handler.Filename)
	if err != nil {
		h.log.Error("Ошибка сохранения файла", err)
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}

	h.log.Info(fmt.Sprintf("\u0424\u0430\u0439\u043b \u0443\u0441\u043f\u0435\u0448\u043d\u043e \u0437\u0430\u0433\u0440\u0443\u0436\u0435\u043d: %s", filePath))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"file_path": "%s"}`, filePath)))
}

// DownloadFileHandler обрабатывает запрос на скачивание файла с сервера.
// @Summary Скачать файл
// @Description Позволяет скачать файл, загруженный на сервер по его имени.
// @Tags файлы
// @Param filename path string true "Имя файла для скачивания"
// @Produce octet-stream
// @Success 200 {file} file "Файл для скачивания"
// @Failure 404 {object} string "Файл не найден"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /admin/files/{filename} [get]
func (h *FileHandler) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на скачивание файла")

	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath, err := h.service.GetFilePath(filename)
	if err != nil {
		h.log.Error("Файл не найден", err)
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		h.log.Error("Не удалось открыть файл", err)
		http.Error(w, "Не удалось открыть файл", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	h.log.Info(fmt.Sprintf("Определённый MIME-тип для файла %s: %s", filename, mimeType))

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", mimeType)

	if _, err := io.Copy(w, file); err != nil {
		h.log.Error("Ошибка при отправке файла", err)
		http.Error(w, "Ошибка при отправке файла", http.StatusInternalServerError)
	}
}
