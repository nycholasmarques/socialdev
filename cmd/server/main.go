package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	configs "github.com/nycholasmarques/socialdev/config"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/models"
	"github.com/nycholasmarques/socialdev/internal/routes"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := configs.LoadConfig()

	// Conectar ao banco de dados com GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // Desativa constraints automáticas
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Criar tabelas na ordem correta, sem constraints
	err = gormDB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	// Adicionar constraints de chave estrangeira manualmente
	err = gormDB.Exec(`
		ALTER TABLE posts 
		ADD CONSTRAINT fk_posts_user 
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE;

		ALTER TABLE comments 
		ADD CONSTRAINT fk_comments_user 
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		ADD CONSTRAINT fk_comments_post 
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE;

		ALTER TABLE likes 
		ADD CONSTRAINT fk_likes_user 
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		ADD CONSTRAINT fk_likes_post 
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE;
	`).Error
	if err != nil {
		log.Fatalf("Failed to add foreign key constraints: %v", err)
	}

	// Manter compatibilidade com sqlc temporariamente
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Configurar Gin
	r := gin.Default()

	// Configuração CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"application/json", "Content-Type"},
	})
	r.Use(gin.WrapH(corsMiddleware.Handler(r))) // Corrigido para usar diretamente o handler

	// Configurar rotas
	routes.SetupRoutes(r, &db.Queries{}, gormDB)

	// Iniciar o servidor
	fmt.Printf("Server working on http://localhost:%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}