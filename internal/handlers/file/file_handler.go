package file

import (
	"ROOmail/pkg/logger"
	"fmt"
	"net/http"
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
// @Router /admin/file/upload [post]
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
