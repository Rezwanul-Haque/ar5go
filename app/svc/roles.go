package svc

import (
	"ar5go/app/domain"
	"ar5go/app/serializers"
	"ar5go/infra/errors"
)

type IRoles interface {
	CreateRole(req *serializers.RoleReq) (*domain.Role, *errors.RestErr)
	GetRole(id uint) (*serializers.RoleResp, *errors.RestErr)
	UpdateRole(roleID uint, req serializers.RoleReq) *errors.RestErr
	DeleteRole(id uint) *errors.RestErr
	ListRoles() ([]*serializers.RoleResp, *errors.RestErr)
	SetRolePermissions(req *serializers.RolePermissionsReq) *errors.RestErr
	GetRolePermissions(roleID int) ([]*domain.Permission, *errors.RestErr)
}
