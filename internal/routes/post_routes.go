package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupPostRoutes(r *http.ServeMux, queries *db.Queries) {
	postRepo := repository.NewPostRepository(queries)
	postService := services.NewPostSerivice(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	r.HandleFunc("POST /posts", middlewares.ProtectedRoutes(postHandler.CreatePostHandler))
	r.HandleFunc("PATCH /posts/{post_id}", middlewares.ProtectedRoutes(postHandler.UpdatePostHandler))
	r.HandleFunc("GET /posts/{post_id}", middlewares.ProtectedRoutes(postHandler.GetPostHandler))
	r.HandleFunc("GET /posts", middlewares.ProtectedRoutes(postHandler.GetPostsHandler))
	r.HandleFunc("GET /users/{user_id}/posts", middlewares.ProtectedRoutes(postHandler.GetUserPostsHandler))
	r.HandleFunc("DELETE /posts/{post_id}", middlewares.ProtectedRoutes(postHandler.DeletePostHandler))
}
