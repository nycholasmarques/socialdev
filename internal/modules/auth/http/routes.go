// internal/modules/auth/delivery/http/routes.go
package http

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/modules/auth/repository"
	"github.com/nycholasmarques/socialdev/internal/modules/auth/service"
)

func SetupAuthRoutes(mux *http.ServeMux, queries *db.Queries) {
	
	authRepository := repository.NewAuthRepository(queries)
	authService := services.NewAuthService(authRepository)
	handler := NewAuthHandler(*authService)
	oAuthHandler := NewOOAuthHandler(*authService)
	
	mux.HandleFunc("POST /auth/login", handler.LoginAuthHandler)
	mux.HandleFunc("GET /auth/google", oAuthHandler.HandleGoogleLogin)
	mux.HandleFunc("GET /auth/google/callback", oAuthHandler.HandleGoogleCallback)
	mux.HandleFunc("GET /auth/verify", middlewares.ProtectedRoutes(handler.VerifyTokenInSession))
}