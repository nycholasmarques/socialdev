package post_handler

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	post_repository "github.com/nycholasmarques/socialdev/internal/modules/posts/repository"
	post_services "github.com/nycholasmarques/socialdev/internal/modules/posts/service"
)

func SetupPostRoutes(r *http.ServeMux, queries *db.Queries) {
	
	repo := post_repository.NewPostRepository(queries)
	services := post_services.NewPostSerivice(repo)
	handler := NewPostHandler(services)
	
	r.HandleFunc("POST /posts", middlewares.ProtectedRoutes(handler.CreatePostHandler))
	r.HandleFunc("PATCH /posts/{post_id}", middlewares.ProtectedRoutes(handler.UpdatePostHandler))
	r.HandleFunc("GET /posts/{post_id}", middlewares.ProtectedRoutes(handler.GetPostHandler))
	r.HandleFunc("GET /posts", middlewares.ProtectedRoutes(handler.GetPostsHandler))
	r.HandleFunc("GET /users/{user_id}/posts", middlewares.ProtectedRoutes(handler.GetUserPostsHandler))
	r.HandleFunc("DELETE /posts/{post_id}", middlewares.ProtectedRoutes(handler.DeletePostHandler))
}