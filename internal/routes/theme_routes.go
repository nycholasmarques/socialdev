package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupThemeRoutes(r *http.ServeMux, queries *db.Queries) {
	themeRepo := repository.NewThemeRepository(queries)
	themeService := services.NewThemeService(themeRepo)
	themeHandler := handlers.NewThemeHandler(themeService)

	r.HandleFunc("POST /themes", themeHandler.CreateThemeHandler)
	r.HandleFunc("GET /themes/{theme_id}", themeHandler.GetThemeHandler)
	r.HandleFunc("GET /themes", themeHandler.ListAllThemesHandler)
	r.HandleFunc("PATCH /themes/{theme_id}", themeHandler.UpdateThemeHandler)
	r.HandleFunc("DELETE /themes/{theme_id}", themeHandler.DeleteThemeHandler)

	r.HandleFunc("POST /users/themes/{theme_id}", middlewares.ProtectedRoutes(themeHandler.AddUserThemeHandler))
	r.HandleFunc("GET /users/themes", middlewares.ProtectedRoutes(themeHandler.ListUserThemesHandler))
	r.HandleFunc("GET /users/themes/active", middlewares.ProtectedRoutes(themeHandler.UserHasActiveThemeHandler))
	r.HandleFunc("PATCH /users/themes/{theme_id}", middlewares.ProtectedRoutes(themeHandler.UpdateUserThemeHandler))
	r.HandleFunc("DELETE /users/{user_id}/themes/{theme_id}", themeHandler.RemoveUserThemeHandler)
}
