package file_handlers

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
)

func SetupFileRoutes(r *http.ServeMux, queries *db.Queries) {
	fileHandler := NewFileHandler()

	r.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		avatarPath, err := fileHandler.UploadFile(r)
		if err != nil {
			http.Error(w, "Erro ao fazer upload do arquivo", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"filePath": "` + avatarPath + `"}`))
	})
}