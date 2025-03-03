package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/db"
)

type RoleRepository struct {
	queries *db.Queries
}

func NewRoleRepository(queries *db.Queries) *RoleRepository {
	return &RoleRepository{queries: queries}
}

func GetRoleRepository(queries *db.Queries) *RoleRepository {
	return &RoleRepository{queries: queries}
}

func (r *RoleRepository) CreateRole(ctx context.Context, params db.CreateRoleParams) (db.Role, error) {
	role, err := r.queries.CreateRole(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return role, err
}

func (r *RoleRepository) GetRole(ctx context.Context, role_id uuid.UUID) (db.Role, error) {
	role, err := r.queries.GetRole(ctx, role_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return role, err
}

func (r *RoleRepository) GetAllRoles(ctx context.Context) ([]db.Role, error) {
	allRoles, err := r.queries.GetAllRoles(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
	return allRoles, err
}

func (r *RoleRepository) GetRoleWithName(ctx context.Context, name string) ([]db.Role, error) {
	usernamePattern := name + "%"
	getRoleWithUsername, err := r.queries.GetRoleWithName(ctx, usernamePattern)
	if err != nil {
		log.Printf("%v", err)
	}
	return getRoleWithUsername, err
}

func (r *RoleRepository) UpdateRole(ctx context.Context, params db.UpdateRoleParams) (db.Role, error) {
	role, err := r.queries.UpdateRole(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return role, err
}

func (r *RoleRepository) DeleteRole(ctx context.Context, role_id uuid.UUID) (error) {
	err := r.queries.DeleteRole(ctx, role_id)
	if err != nil {
		log.Printf("%v", err)
	}
	return err
}

func (r *RoleRepository) CreateUserWithRole(ctx context.Context, params db.CreateUserWithRolesParams) (db.UserRole, error) {
	user, err := r.queries.CreateUserWithRoles(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return user, err
}

func (r *RoleRepository) DeleteUserWithRole(ctx context.Context, params db.DeleteUserWithRolesParams) (error) {
	err := r.queries.DeleteUserWithRoles(ctx, params)
	if err != nil {
		log.Printf("%v", err)
	}
	return err
}