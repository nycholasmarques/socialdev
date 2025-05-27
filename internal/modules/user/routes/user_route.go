package user_routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	user_handler "github.com/nycholasmarques/socialdev/internal/modules/user/handler"
	user_repository "github.com/nycholasmarques/socialdev/internal/modules/user/repository"
	user_service "github.com/nycholasmarques/socialdev/internal/modules/user/service"
)

func SetupUserRoutes(r *http.ServeMux, queries *db.Queries) {
	
	repo := user_repository.NewUserRepository(queries)
	services := user_service.NewUserService(repo)
	handler := user_handler.NewUserHandler(services)
	
	r.HandleFunc("POST /users", handler.CreateUserHandler)
	r.HandleFunc("GET /users/{user_id}", handler.GetUserHandler)
	r.HandleFunc("GET /users", middlewares.ProtectedRoutes(handler.GetAllUsersHandler))
	r.HandleFunc("PATCH /users", middlewares.ProtectedRoutes(handler.UpdateUserHandler))
	r.HandleFunc("DELETE /users/{user_id}", middlewares.ProtectedRoutes(handler.DeleteUserHandler))
}