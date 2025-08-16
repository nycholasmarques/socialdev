package main

import (
	"fmt"
	"log"
	"net/http"


	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	configs "github.com/nycholasmarques/socialdev/config"
	"github.com/nycholasmarques/socialdev/internal/db"
	auth_routes "github.com/nycholasmarques/socialdev/internal/modules/auth/routes"
	comment_routes "github.com/nycholasmarques/socialdev/internal/modules/comment/routes"
	favorite_routes "github.com/nycholasmarques/socialdev/internal/modules/favorite/routes"
	file_routes "github.com/nycholasmarques/socialdev/internal/modules/file/routes"
	like_routes "github.com/nycholasmarques/socialdev/internal/modules/like/routes"
	post_route "github.com/nycholasmarques/socialdev/internal/modules/posts/routes"
	role_routes "github.com/nycholasmarques/socialdev/internal/modules/role/routes"
	theme_routes "github.com/nycholasmarques/socialdev/internal/modules/theme/routes"
	user_routes "github.com/nycholasmarques/socialdev/internal/modules/user/routes"
)

func main() {
	cfg := configs.LoadConfig()

	conn, err := configs.NewConnectionDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	queries := db.New(conn)
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		auth_routes.SetupAuthRoutes(r, queries)
		user_routes.SetupUserRoutes(r, queries)
		role_routes.SetupRoleRoutes(r, queries)
		theme_routes.SetupThemeRoutes(r, queries)
		post_route.SetupPostRoutes(r, queries)
		like_routes.SetupLikeRoutes(r, queries)
		comment_routes.SetupCommentRoutes(r, queries)
		favorite_routes.SetupFavoriteRoutes(r, queries)
		file_routes.SetupFileRoutes(r, queries)
	})

	fmt.Println("Server working on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}