package like_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	like_handler "github.com/nycholasmarques/socialdev/internal/modules/like/handler"
	like_repository "github.com/nycholasmarques/socialdev/internal/modules/like/repository"
	like_service "github.com/nycholasmarques/socialdev/internal/modules/like/service"
)

func SetupLikeRoutes(r chi.Router, queries *db.Queries) {
	repo := like_repository.NewLikeRepository(queries)
	services := like_service.NewLikeService(repo)
	handler := like_handler.NewLikeHandler(services)

	r.With(middlewares.ProtectedRoutes).Post("/likes/posts/{post_id}", handler.CreateLikeHandler)
	r.With(middlewares.ProtectedRoutes).Delete("/likes/posts/{post_id}", handler.DeleteLikeHandler)
	r.With(middlewares.ProtectedRoutes).Get("/likes/posts/{post_id}", handler.GetLikeHandler)
}
