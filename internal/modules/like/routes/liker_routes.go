package like_routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	like_handler "github.com/nycholasmarques/socialdev/internal/modules/like/handler"
	like_repository "github.com/nycholasmarques/socialdev/internal/modules/like/repository"
	like_service "github.com/nycholasmarques/socialdev/internal/modules/like/service"
)

func SetupLikeRoutes(r *http.ServeMux, queries *db.Queries) {
	
	repo := like_repository.NewLikeRepository(queries)
	services := like_service.NewLikeService(repo)
	handler := like_handler.NewLikeHandler(services)
	
	
	r.HandleFunc("POST /likes/posts/{post_id}", middlewares.ProtectedRoutes(handler.CreateLikeHandler))
	r.HandleFunc("DELETE /likes/posts/{post_id}", middlewares.ProtectedRoutes(handler.DeleteLikeHandler))
	r.HandleFunc("GET /likes/posts/{post_id}", middlewares.ProtectedRoutes(handler.GetLikeHandler))
}