package routes

import (

	"github.com/gin-gonic/gin"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
	"gorm.io/gorm"
)

func SetupPostRoutes(r *gin.RouterGroup, queries *db.Queries, gormDB *gorm.DB) {
	postRepo := repository.NewPostRepository(queries)
	postService := services.NewPostSerivice(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	r.POST("/posts", middlewares.ProtectedRoutes(), postHandler.CreatePostHandler)
	// r.PATCH("/posts/:post_id", middlewares.ProtectedRoutes(), postHandler.UpdatePostHandler)
	// r.GET("/posts/:post_id", postHandler.GetPostHandler)
	// r.GET("/posts", postHandler.GetPostsHandler)
	// r.GET("/users/:user_id/posts", middlewares.ProtectedRoutes(), postHandler.GetUserPostsHandler)
	// r.DELETE("/posts/:post_id", middlewares.ProtectedRoutes(), postHandler.DeletePostHandler)
}