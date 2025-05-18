package file_routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	file_handlers "github.com/nycholasmarques/socialdev/internal/modules/file/handler"
)

func SetupFileRoutes(r *http.ServeMux, queries *db.Queries) {
	fileHandler := file_handlers.NewFileHandler()

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