package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupLikeRoutes(r *http.ServeMux, queries *db.Queries) {
	likeRepo := repository.NewLikeRepository(queries)
	likeService := services.NewLikeService(likeRepo)
	likeHandler := handlers.NewLikeHandler(likeService)

	r.HandleFunc("POST /likes/users/{user_id}/posts/{post_id}", likeHandler.CreateLikeHandler)
	r.HandleFunc("DELETE /likes/users/{user_id}/posts/{post_id}", likeHandler.DeleteLikeHandler)
}