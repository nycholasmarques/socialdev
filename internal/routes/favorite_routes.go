package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupFavoriteRoutes(r *http.ServeMux, queries *db.Queries) {
	favoriteRepo := repository.NewFavoriteRepository(queries)
	favoriteService := services.NewFavoriteService(favoriteRepo)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteService)

	r.HandleFunc("POST /favorites/{post_id}", middlewares.ProtectedRoutes(favoriteHandler.AddFavoriteHandler))
	r.HandleFunc("DELETE /favorites/{post_id}", middlewares.ProtectedRoutes(favoriteHandler.RemoveFavoriteHandler))
	r.HandleFunc("GET /favorites", middlewares.ProtectedRoutes(favoriteHandler.GetUserFavoritesHandler))
}
