package role_routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nycholasmarques/socialdev/internal/db"
	role_handler "github.com/nycholasmarques/socialdev/internal/modules/role/handler"
	role_repository "github.com/nycholasmarques/socialdev/internal/modules/role/repository"
	role_service "github.com/nycholasmarques/socialdev/internal/modules/role/service"
)

func SetupRoleRoutes(r chi.Router, queries *db.Queries) {
	repo := role_repository.NewRoleRepository(queries)
	services := role_service.NewRoleService(repo)
	handler := role_handler.NewRoleHandler(services)

	r.Post("/roles", handler.CreateRoleHandler)
	r.Get("/roles/{role_id}", handler.GetRoleHandler)
	r.Get("/roles", handler.GetAllRolesHandler)
	r.Patch("/roles/{role_id}", handler.UpdateRoleHandler)
	r.Delete("/roles/{role_id}", handler.DeleteRoleHandler)

	r.Post("/users/{user_id}/roles/{role_id}", handler.CreateUserWithRoleHandler)
	r.Delete("/users/{user_id}/roles/{role_id}", handler.DeleteUserWithRoleHandler)
}
