package file

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"path/filepath"
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
	filePath := filepath.Join(s.uploadDir, filename)

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