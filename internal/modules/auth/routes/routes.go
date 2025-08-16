package auth_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	auth_handlers "github.com/nycholasmarques/socialdev/internal/modules/auth/http"
	"github.com/nycholasmarques/socialdev/internal/modules/auth/repository"
	auth_service "github.com/nycholasmarques/socialdev/internal/modules/auth/service"
)

func SetupAuthRoutes(r chi.Router, queries *db.Queries) {
	authRepository := repository.NewAuthRepository(queries)
	authService := auth_service.NewAuthService(authRepository)
	handler := auth_handlers.NewAuthHandler(*authService)
	oAuthHandler := auth_handlers.NewOOAuthHandler(*authService)

	r.Post("/auth/login", handler.LoginAuthHandler)
	r.Get("/auth/google", oAuthHandler.HandleGoogleLogin)
	r.Get("/auth/google/callback", oAuthHandler.HandleGoogleCallback)
	r.With(middlewares.ProtectedRoutes).Get("/auth/verify", handler.VerifyTokenInSession)
}
