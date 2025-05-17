package comment_handlers

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	comment_repository "github.com/nycholasmarques/socialdev/internal/modules/comment/repository"
	comment_service "github.com/nycholasmarques/socialdev/internal/modules/comment/service"
)

func SetupCommentRoutes(r *http.ServeMux, queries *db.Queries) {
	
	repo := comment_repository.NewCommentRepository(queries)
	services := comment_service.NewCommentService(repo)
	handler := NewCommentHandler(services)
	
	r.HandleFunc("POST /comments/posts/{post_id}", middlewares.ProtectedRoutes(handler.CreateCommentHandler))
	r.HandleFunc("GET /comments/posts/{post_id}", middlewares.ProtectedRoutes(handler.GetAllCommentsHandler))
	r.HandleFunc("GET /comments/{comment_id}", middlewares.ProtectedRoutes(handler.GetCommentHandler))
	r.HandleFunc("PATCH /comments/{comment_id}", middlewares.ProtectedRoutes(handler.UpdateCommentHandler))
	r.HandleFunc("DELETE /comments/{comment_id}", middlewares.ProtectedRoutes(handler.DeleteCommentHandler))
}