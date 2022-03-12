package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/infra/conn/db"
	"ar5go/infra/errors"
	"context"
)

type roles struct {
	ctx context.Context
	DB  db.DatabaseClient
}

// NewRolesRepository will create an object that represent the Roles.Repository implementations
func NewRolesRepository(ctx context.Context, dbc db.DatabaseClient) repository.IRoles {
	return &roles{
		ctx: ctx,
		DB:  dbc,
	}
}

func (r *roles) CreateRole(role *domain.Role) (*domain.Role, *errors.RestErr) {
	return r.DB.CreateRole(role)
}

func (r *roles) GetRole(roleID uint) (*domain.Role, *errors.RestErr) {
	return r.DB.GetRole(roleID)
}

func (r *roles) UpdateRole(role *domain.Role) *errors.RestErr {
	return r.DB.UpdateRole(role)
}

func (r *roles) RemoveRole(id uint) *errors.RestErr {
	return r.DB.RemoveRole(id)
}

func (r *roles) ListRoles() ([]*domain.Role, *errors.RestErr) {
	return r.DB.ListRoles()
}

func (r *roles) SetRolePermissions(rolePerms *domain.RolePermissions) *errors.RestErr {
	return r.DB.SetRolePermissions(rolePerms)
}

func (r *roles) GetRolePermissions(roleID int) ([]*domain.Permission, *errors.RestErr) {
	return r.DB.GetRolePermissions(roleID)
}
