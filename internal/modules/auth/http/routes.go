// internal/modules/auth/delivery/http/routes.go
package http

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/modules/auth/service"
)

func SetupAuthRoutes(mux *http.ServeMux, authService services.AuthService) {
	
	handler := NewAuthHandler(authService)
	oAuthHandler := NewOOAuthHandler(authService)
	
	mux.HandleFunc("POST /auth/login", handler.LoginAuthHandler)
	mux.HandleFunc("GET /auth/google", oAuthHandler.HandleGoogleLogin)
	mux.HandleFunc("GET /auth/google/callback", oAuthHandler.HandleGoogleCallback)
	mux.HandleFunc("GET /auth/verify", middlewares.ProtectedRoutes(handler.VerifyTokenInSession))
}