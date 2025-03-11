package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupPostRoutes(r *http.ServeMux, queries *db.Queries) {
	postRepo := repository.NewPostRepository(queries)
	postService := services.NewPostSerivice(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	r.HandleFunc("POST /posts", postHandler.CreatePostHandler)
}