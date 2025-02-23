package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
)

func SetupRoutes(queries *db.Queries) *http.ServeMux {
	r := http.NewServeMux()

	// Criar um grupo de rotas com o prefixo /api
	apiRouter := http.NewServeMux()
	SetupAPIRoutes(apiRouter, queries)

	// Registrar o grupo de rotas /api no ServeMux principal
	r.Handle("/api/", http.StripPrefix("/api", apiRouter))

	return r
}