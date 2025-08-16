package user_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	user_handler "github.com/nycholasmarques/socialdev/internal/modules/user/handler"
	user_repository "github.com/nycholasmarques/socialdev/internal/modules/user/repository"
	user_service "github.com/nycholasmarques/socialdev/internal/modules/user/service"
)

func SetupUserRoutes(r chi.Router, queries *db.Queries) {
	repo := user_repository.NewUserRepository(queries)
	services := user_service.NewUserService(repo)
	handler := user_handler.NewUserHandler(services)

	r.Post("/users", handler.CreateUserHandler)
	r.Get("/users/{user_id}", handler.GetUserHandler)
	r.With(middlewares.ProtectedRoutes).Get("/users", handler.GetAllUsersHandler)
	r.With(middlewares.ProtectedRoutes).Patch("/users/{user_id}", handler.UpdateUserHandler)
	r.With(middlewares.ProtectedRoutes).Delete("/users/{user_id}", handler.DeleteUserHandler)
}
