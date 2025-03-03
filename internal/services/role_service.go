package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
	"github.com/nycholasmarques/socialdev/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (r *RoleService) CreateRole(role_id uuid.UUID, name string) (db.Role, error) {
	ctx := context.Background()

	params := db.CreateRoleParams{RoleID: role_id, Name: name}
	return r.roleRepo.CreateRole(ctx, params)
}

func (r *RoleService) GetRole(role_id uuid.UUID) (db.Role, error) {
	ctx := context.Background()
	return r.roleRepo.GetRole(ctx, role_id)
}

func (r *RoleService) GetAllRoles() ([]db.Role, error) {
	ctx := context.Background()
	return r.roleRepo.GetAllRoles(ctx)
}

func (r *RoleService) GetRoleWithName(name string) ([]db.Role, error) {
	ctx := context.Background()
	return r.roleRepo.GetRoleWithName(ctx, name)
}

func (r *RoleService) UpdateRole(role_id uuid.UUID, name string) (db.Role, error) {
	ctx := context.Background()

	params := db.UpdateRoleParams{RoleID: role_id, Name: name}
	return r.roleRepo.UpdateRole(ctx, params)
}

func (r *RoleService) DeleteRole(role_id uuid.UUID) (error) {
	ctx := context.Background()
	return r.roleRepo.DeleteRole(ctx, role_id)
}

func (r *RoleService) CreateUserWithRole(user_id uuid.UUID, role_id uuid.UUID) (db.UserRole, error) {
	ctx := context.Background()

	params := db.CreateUserWithRolesParams{UserID: user_id, RoleID: role_id}
	return r.roleRepo.CreateUserWithRole(ctx, params)
}

func (r *RoleService) DeleteUserWithRole(user_id uuid.UUID, role_id uuid.UUID) (error) {
	ctx := context.Background()

	params := db.DeleteUserWithRolesParams{UserID: user_id, RoleID: role_id}
	return r.roleRepo.DeleteUserWithRole(ctx, params)
}