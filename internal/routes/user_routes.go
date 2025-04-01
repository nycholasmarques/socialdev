package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupUserRoutes(r *http.ServeMux, queries *db.Queries) {
	userRepo := repository.NewUserRepository(queries)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r.HandleFunc("POST /users", userHandler.CreateUserHandler)
	r.HandleFunc("GET /users/{user_id}", userHandler.GetUserHandler)
	r.HandleFunc("GET /users", userHandler.GetAllUsersHandler)
	r.HandleFunc("PATCH /users/{user_id}", middlewares.ProtectedRoutes(userHandler.UpdateUserHandler))
	r.HandleFunc("DELETE /users/{user_id}", middlewares.ProtectedRoutes(userHandler.DeleteUserHandler))

	r.HandleFunc("POST /auth/login", userHandler.LoginUserHandler)
	r.HandleFunc("GET /auth/google", userHandler.HandleGoogleLogin)
	r.HandleFunc("GET /auth/google/callback", userHandler.HandleGoogleCallback)
	r.HandleFunc("GET /auth/verify", middlewares.ProtectedRoutes(userHandler.VerifyTokenInSession))
}
