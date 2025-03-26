 package routes

// import (
// 	"net/http"

// 	"github.com/nycholasmarques/socialdev/internal/db"
// 	"github.com/nycholasmarques/socialdev/internal/handlers"
// 	"github.com/nycholasmarques/socialdev/internal/middlewares"
// 	"github.com/nycholasmarques/socialdev/internal/repository"
// 	"github.com/nycholasmarques/socialdev/internal/services"
// )

// func SetupCommentRoutes(r *http.ServeMux, queries *db.Queries) {
// 	commentRepo := repository.NewCommentRepository(queries)
// 	commentService := services.NewCommentService(commentRepo)
// 	commentHandler := handlers.NewCommentHandler(commentService)

// 	r.HandleFunc("POST /comments/posts/{post_id}", middlewares.ProtectedRoutes(commentHandler.CreateCommentHandler))
// 	r.HandleFunc("GET /comments/posts/{post_id}", middlewares.ProtectedRoutes(commentHandler.GetAllCommentsHandler))
// 	r.HandleFunc("GET /comments/{comment_id}", middlewares.ProtectedRoutes(commentHandler.GetCommentHandler))
// 	r.HandleFunc("PATCH /comments/{comment_id}", middlewares.ProtectedRoutes(commentHandler.UpdateCommentHandler))
// 	r.HandleFunc("DELETE /comments/{comment_id}", middlewares.ProtectedRoutes(commentHandler.DeleteCommentHandler))
// }
