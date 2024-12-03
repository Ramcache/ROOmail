package file

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileInterface interface {
	SaveFile(file io.Reader, filename string) (string, error)
}

type FileService struct {
	uploadDir string
	db        *pgxpool.Pool
}

func NewFileService(uploadDir string, db *pgxpool.Pool) *FileService {
	return &FileService{
		uploadDir: uploadDir,
		db:        db,
	}
}

func (s *FileService) SaveFile(file io.Reader, filename string) (string, error) {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Unix()
	newFilename := fmt.Sprintf("%d%s", timestamp, ext)
	filePath := filepath.Join(s.uploadDir, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filePath, nil
}

func (s *FileService) GetFilePath(filename string) (string, error) {
	filePath := filepath.Join(s.uploadDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("Файл не найден: %s", filename)
	}

	return filePath, nil
}
