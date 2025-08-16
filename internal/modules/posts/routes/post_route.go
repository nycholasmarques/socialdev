package post_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	post_handler "github.com/nycholasmarques/socialdev/internal/modules/posts/handler"
	post_repository "github.com/nycholasmarques/socialdev/internal/modules/posts/repository"
	post_services "github.com/nycholasmarques/socialdev/internal/modules/posts/service"
)

func SetupPostRoutes(r chi.Router, queries *db.Queries) {
	repo := post_repository.NewPostRepository(queries)
	services := post_services.NewPostSerivice(repo)
	handler := post_handler.NewPostHandler(services)

	r.With(middlewares.ProtectedRoutes).Post("/posts", handler.CreatePostHandler)
	r.With(middlewares.ProtectedRoutes).Patch("/posts/{post_id}", handler.UpdatePostHandler)
	r.With(middlewares.ProtectedRoutes).Get("/posts/{post_id}", handler.GetPostHandler)
	r.With(middlewares.ProtectedRoutes).Get("/posts", handler.GetPostsHandler)
	r.With(middlewares.ProtectedRoutes).Get("/users/{user_id}/posts", handler.GetUserPostsHandler)
	r.With(middlewares.ProtectedRoutes).Delete("/posts/{post_id}", handler.DeletePostHandler)
}
