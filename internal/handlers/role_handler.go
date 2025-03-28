package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/services"
)

type RoleHandler struct {
	roleService *services.RoleService
}

type RoleWithUser struct {
	RoleID uuid.UUID `json:"role_id"`
	UserID uuid.UUID `json:"user_id"`
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	rb := r.Body
	reader, err := io.ReadAll(rb)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	rb.Close()

	var roleCreate *db.CreateRoleParams
	err = json.Unmarshal(reader, &roleCreate)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	role_id := uuid.New()

	role, err := h.roleService.CreateRole(role_id, roleCreate.Name)
	if err != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(role)
	if err != nil {
		http.Error(w, "Failed to marshal role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *RoleHandler) GetRoleHandler(w http.ResponseWriter, r *http.Request) {
	role_id := r.URL.Path
	roleID := strings.TrimPrefix(role_id, "/roles/")

	uuid, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	role, err := h.roleService.GetRole(uuid)
	if err != nil {
		http.Error(w, "Failed to get role", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(role)
	if err != nil {
		http.Error(w, "Failed to marshal role", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *RoleHandler) GetAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		http.Error(w, "Failed to get roles", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(roles)
	if err != nil {
		http.Error(w, "Failed to marshal roles", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *RoleHandler) GetRoleWithNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Println(name)
	roles, err := h.roleService.GetRoleWithName(name)
	if err != nil {
		http.Error(w, "Failed to get role", http.StatusInternalServerError)
		return
	}

	if len(roles) == 0 {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	res, err := json.Marshal(roles)
	if err != nil {
		http.Error(w, "Failed to marshal role", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *RoleHandler) UpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	role_id := r.URL.Path
	roleID := strings.TrimPrefix(role_id, "/roles/")

	uuid, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	role, err := h.roleService.GetRole(uuid)
	if err != nil {
		http.Error(w, "Role does not exist", http.StatusNotFound)
		return
	}

	rb := r.Body
	reader, err := io.ReadAll(rb)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	rb.Close()

	var roleUpdate *db.UpdateRoleParams
	err = json.Unmarshal(reader, &roleUpdate)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	role, err = h.roleService.UpdateRole(role.RoleID, roleUpdate.Name)
	if err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(role)
	if err != nil {
		http.Error(w, "Failed to marshal role", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (h *RoleHandler) DeleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	role_id := r.URL.Path
	roleID := strings.TrimPrefix(role_id, "/roles/")

	uuid, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	_, err = h.roleService.GetRole(uuid)
	if err != nil {
		http.Error(w, "Role does not exist", http.StatusNotFound)
		return
	}

	err = h.roleService.DeleteRole(uuid)
	if err != nil {
		http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Role deleted successfully!"))
}

func (h *RoleHandler) CreateUserWithRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros da URL
	userID := r.PathValue("user_id")
	roleID := r.PathValue("role_id")

	if userID == "" || roleID == "" {
		http.Error(w, "Missing user_id or role_id in URL", http.StatusBadRequest)
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	roleIdUUID, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	user, err := h.roleService.CreateUserWithRole(userIdUUID, roleIdUUID)
	if err != nil {
		http.Error(w, "Failed to create user with role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (h *RoleHandler) DeleteUserWithRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros da URL
	userID := r.PathValue("user_id")
	roleID := r.PathValue("role_id")

	if userID == "" || roleID == "" {
		http.Error(w, "Missing user_id or role_id in URL", http.StatusBadRequest)
		return
	}

	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	roleIdUUID, err := uuid.Parse(roleID)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.roleService.DeleteUserWithRole(userIdUUID, roleIdUUID); err != nil {
		http.Error(w, "Failed to delete user with role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User role removed successfully"))
}
