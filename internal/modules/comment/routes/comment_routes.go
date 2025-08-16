package comment_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	comment_handlers "github.com/nycholasmarques/socialdev/internal/modules/comment/handler"
	comment_repository "github.com/nycholasmarques/socialdev/internal/modules/comment/repository"
	comment_service "github.com/nycholasmarques/socialdev/internal/modules/comment/service"
)

func SetupCommentRoutes(r chi.Router, queries *db.Queries) {
	repo := comment_repository.NewCommentRepository(queries)
	services := comment_service.NewCommentService(repo)
	handler := comment_handlers.NewCommentHandler(services)

	r.With(middlewares.ProtectedRoutes).Post("/comments/posts/{post_id}", handler.CreateCommentHandler)
	r.With(middlewares.ProtectedRoutes).Get("/comments/posts/{post_id}", handler.GetAllCommentsHandler)
	r.With(middlewares.ProtectedRoutes).Get("/comments/{comment_id}", handler.GetCommentHandler)
	r.With(middlewares.ProtectedRoutes).Patch("/comments/{comment_id}", handler.UpdateCommentHandler)
	r.With(middlewares.ProtectedRoutes).Delete("/comments/{comment_id}", handler.DeleteCommentHandler)
}
