package role_routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	role_handler "github.com/nycholasmarques/socialdev/internal/modules/role/http"
	role_repository "github.com/nycholasmarques/socialdev/internal/modules/role/repository"
	role_service "github.com/nycholasmarques/socialdev/internal/modules/role/service"
)

func SetupUserRoutes(r *http.ServeMux, queries *db.Queries) {
	repo := role_repository.NewRoleRepository(queries)
	services := role_service.NewRoleService(repo)
	handler := role_handler.NewRoleHandler(services)

	r.HandleFunc("POST /roles", handler.CreateRoleHandler)
	r.HandleFunc("GET /roles/{role_id}", handler.GetRoleHandler)
	r.HandleFunc("GET /roles", handler.GetAllRolesHandler)
	r.HandleFunc("PATCH /roles/{role_id}", handler.UpdateRoleHandler)
	r.HandleFunc("DELETE /roles/{role_id}", handler.DeleteRoleHandler)

	r.HandleFunc("POST /users/{user_id}/roles/{role_id}", handler.CreateUserWithRoleHandler)
	r.HandleFunc("DELETE /users/{user_id}/roles/{role_id}", handler.DeleteUserWithRoleHandler)
}