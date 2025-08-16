package theme_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	theme_handler "github.com/nycholasmarques/socialdev/internal/modules/theme/handler"
	theme_repository "github.com/nycholasmarques/socialdev/internal/modules/theme/repository"
	theme_service "github.com/nycholasmarques/socialdev/internal/modules/theme/service"
)

func SetupThemeRoutes(r chi.Router, queries *db.Queries) {
	themeRepo := theme_repository.NewThemeRepository(queries)
	themeService := theme_service.NewThemeService(themeRepo)
	themeHandler := theme_handler.NewThemeHandler(themeService)

	r.Post("/themes", themeHandler.CreateThemeHandler)
	r.Get("/themes/{theme_id}", themeHandler.GetThemeHandler)
	r.Get("/themes", themeHandler.ListAllThemesHandler)
	r.Patch("/themes/{theme_id}", themeHandler.UpdateThemeHandler)
	r.Delete("/themes/{theme_id}", themeHandler.DeleteThemeHandler)

	r.With(middlewares.ProtectedRoutes).Post("/users/themes/{theme_id}", themeHandler.AddUserThemeHandler)
	r.With(middlewares.ProtectedRoutes).Get("/users/themes", themeHandler.ListUserThemesHandler)
	r.With(middlewares.ProtectedRoutes).Get("/users/themes/active", themeHandler.UserHasActiveThemeHandler)
	r.With(middlewares.ProtectedRoutes).Patch("/users/themes/{theme_id}", themeHandler.UpdateUserThemeHandler)
	r.Delete("/users/{user_id}/themes/{theme_id}", themeHandler.RemoveUserThemeHandler)
}
