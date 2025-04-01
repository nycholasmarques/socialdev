package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupLikeRoutes(r *http.ServeMux, queries *db.Queries) {
	likeRepo := repository.NewLikeRepository(queries)
	likeService := services.NewLikeService(likeRepo)
	likeHandler := handlers.NewLikeHandler(likeService)

	r.HandleFunc("POST /likes/posts/{post_id}", middlewares.ProtectedRoutes(likeHandler.CreateLikeHandler))
	r.HandleFunc("DELETE /likes/posts/{post_id}", middlewares.ProtectedRoutes(likeHandler.DeleteLikeHandler))
	r.HandleFunc("GET /likes/posts/{post_id}", middlewares.ProtectedRoutes(likeHandler.GetLikeHandler))
}
