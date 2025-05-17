package favorite_handlers

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	favorite_repository "github.com/nycholasmarques/socialdev/internal/modules/favorite/repository"
	favorite_services "github.com/nycholasmarques/socialdev/internal/modules/favorite/service"
)

func SetupFavoriteRoutes(r *http.ServeMux, queries *db.Queries) {
	
	repo := favorite_repository.NewFavoriteRepository(queries)
	services := favorite_services.NewFavoriteService(repo)
	handler := NewFavoriteHandler(services)
	
	r.HandleFunc("POST /favorites/{post_id}", middlewares.ProtectedRoutes(handler.AddFavoriteHandler))
	r.HandleFunc("DELETE /favorites/{post_id}", middlewares.ProtectedRoutes(handler.RemoveFavoriteHandler))
	r.HandleFunc("GET /favorites", middlewares.ProtectedRoutes(handler.GetUserFavoritesHandler))
}