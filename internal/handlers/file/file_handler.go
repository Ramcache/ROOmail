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

func (h *FileHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filePath, err := h.service.SaveFile(file, handler.Filename)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"file_path": "%s"}`, filePath)))
}
