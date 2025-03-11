package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
)

func SetupAPIRoutes(r *http.ServeMux, queries *db.Queries) {
	SetupUserRoutes(r, queries)
	// SetupPostRoutes(r, queries)
}