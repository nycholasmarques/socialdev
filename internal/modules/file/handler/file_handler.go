package file_handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileHandler struct {
	UploadsDir   string
	MaxUploadSize int64
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		UploadsDir:   "./uploads",
		MaxUploadSize: 10 << 20, // 10MB
	}
}

func (h *FileHandler) initUploadsDir() error {
	if _, err := os.Stat(h.UploadsDir); os.IsNotExist(err) {
		return os.MkdirAll(h.UploadsDir, os.ModePerm)
	}
	return nil
}

func (h *FileHandler) UploadFile(r *http.Request) (string, error) {
	if err := h.initUploadsDir(); err != nil {
		return "", err
	}

	r.Body = http.MaxBytesReader(nil, r.Body, h.MaxUploadSize)
	if err := r.ParseMultipartForm(h.MaxUploadSize); err != nil {
		return "", err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}
	file.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		return "", err
	}

	ext := filepath.Ext(header.Filename)
	newFilename := time.Now().Format("20060102150405") + "_" + uuid.New().String() + ext
	filePath := filepath.Join(h.UploadsDir, newFilename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	// Aqui retorna o path relativo para o front usar
	return "/uploads/" + newFilename, nil
}
