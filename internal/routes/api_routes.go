package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
)

func SetupAPIRoutes(r *http.ServeMux, queries *db.Queries) {
	// Registrar as rotas do módulo de usuário
	SetupUserRoutes(r, queries)
}