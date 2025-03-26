package routes

import (

	"github.com/gin-gonic/gin"
	"github.com/nycholasmarques/socialdev/internal/db"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, queries *db.Queries, gormDB *gorm.DB) {
	api := r.Group("/api")


	SetupPostRoutes(api, queries, gormDB)
	// Criar um grupo de rotas com o prefixo /api
	// apiRouter := http.NewServeMux()
	// SetupAPIRoutes(apiRouter, queries)
	// SetupRoleRoutes(apiRouter, queries)
	// SetupPostRoutes(apiRouter, queries)
	// SetupLikeRoutes(apiRouter, queries)
	// SetupCommentRoutes(apiRouter, queries)

	// Registrar o grupo de rotas /api no ServeMux principal
	// r.Handle("/api/", http.StripPrefix("/api", apiRouter))
}