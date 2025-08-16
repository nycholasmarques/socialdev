package favorite_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	favorite_handlers "github.com/nycholasmarques/socialdev/internal/modules/favorite/handler"
	favorite_repository "github.com/nycholasmarques/socialdev/internal/modules/favorite/repository"
	favorite_services "github.com/nycholasmarques/socialdev/internal/modules/favorite/service"
)

func SetupFavoriteRoutes(r chi.Router, queries *db.Queries) {
	repo := favorite_repository.NewFavoriteRepository(queries)
	services := favorite_services.NewFavoriteService(repo)
	handler := favorite_handlers.NewFavoriteHandler(services)

	r.With(middlewares.ProtectedRoutes).Post("/favorites/{post_id}", handler.AddFavoriteHandler)
	r.With(middlewares.ProtectedRoutes).Delete("/favorites/{post_id}", handler.RemoveFavoriteHandler)
	r.With(middlewares.ProtectedRoutes).Get("/favorites", handler.GetUserFavoritesHandler)
}
