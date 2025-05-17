package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/modules/auth/repository"
	services "github.com/nycholasmarques/socialdev/internal/modules/auth/service"
	authHttp "github.com/nycholasmarques/socialdev/internal/modules/auth/http"

)

func SetupRoutes(queries *db.Queries) *http.ServeMux {
	r := http.NewServeMux()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	uploadsDir := filepath.Join(wd, "uploads")
	r.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))



	apiRouter := http.NewServeMux()
	authRepository := repository.NewAuthRepository(queries)
	authService := services.NewAuthService(authRepository)
	authHttp.SetupAuthRoutes(apiRouter, *authService)
	SetupAPIRoutes(apiRouter, queries)
	SetupRoleRoutes(apiRouter, queries)
	SetupPostRoutes(apiRouter, queries)
	SetupLikeRoutes(apiRouter, queries)
	SetupCommentRoutes(apiRouter, queries)
	SetupFavoriteRoutes(apiRouter, queries)
	SetupThemeRoutes(apiRouter, queries)
	SetupFileRoutes(apiRouter, queries)

	r.Handle("/api/", http.StripPrefix("/api", apiRouter))

	return r
}
