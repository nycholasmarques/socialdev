package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupUserRoutes(r *http.ServeMux, queries *db.Queries) {
	// Criar instâncias das camadas
	userRepo := repository.NewUserRepository(queries)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Registrar a rota para criar usuário
	r.HandleFunc("POST /user/create", userHandler.CreateUserHandler)
	r.HandleFunc("GET /user/{user_id}", userHandler.GetUserHandler)
}