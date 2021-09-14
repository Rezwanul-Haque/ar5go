package repository

import (
	"boilerplate/app/domain"
	"boilerplate/infra/errors"
)

type IRoles interface {
	Create(role *domain.Role) (*domain.Role, *errors.RestErr)
	Get(roleID uint) (*domain.Role, *errors.RestErr)
	Update(role *domain.Role) *errors.RestErr
	Remove(id uint) *errors.RestErr
	List() ([]*domain.Role, *errors.RestErr)
	SetRolePermissions(rolePerms *domain.RolePermissions) *errors.RestErr
	GetRolePermissions(roleID int) ([]*domain.Permission, *errors.RestErr)
}
