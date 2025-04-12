package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
)

func SetupRoutes(queries *db.Queries) *http.ServeMux {
	r := http.NewServeMux()

	apiRouter := http.NewServeMux()
	SetupAPIRoutes(apiRouter, queries)
	SetupRoleRoutes(apiRouter, queries)
	SetupPostRoutes(apiRouter, queries)
	SetupLikeRoutes(apiRouter, queries)
	SetupCommentRoutes(apiRouter, queries)
	SetupFavoriteRoutes(apiRouter, queries)
	SetupThemeRoutes(apiRouter, queries)

	r.Handle("/api/", http.StripPrefix("/api", apiRouter))

	return r
}
