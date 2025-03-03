package routes

import (
	"net/http"

	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/handlers"
	"github.com/nycholasmarques/socialdev/internal/repository"
	"github.com/nycholasmarques/socialdev/internal/services"
)

func SetupRoleRoutes(r *http.ServeMux, queries *db.Queries) {
	roleRepo := repository.NewRoleRepository(queries)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := handlers.NewRoleHandler(roleService)

	r.HandleFunc("POST /roles", roleHandler.CreateRoleHandler)
	r.HandleFunc("GET /roles/{role_id}", roleHandler.GetRoleHandler)
	r.HandleFunc("GET /roles", roleHandler.GetAllRolesHandler)
	r.HandleFunc("PATCH /roles/{role_id}", roleHandler.UpdateRoleHandler)
	r.HandleFunc("DELETE /roles/{role_id}", roleHandler.DeleteRoleHandler)

	r.HandleFunc("POST /users/{user_id}/roles/{role_id}", roleHandler.CreateUserWithRoleHandler)
	r.HandleFunc("DELETE /users/{user_id}/roles/{role_id}", roleHandler.DeleteUserWithRoleHandler)	
}