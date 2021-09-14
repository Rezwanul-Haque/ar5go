package svc

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infra/errors"
)

type IRoles interface {
	CreateRole(req *serializers.RoleReq) (*domain.Role, *errors.RestErr)
	GetRule(id uint) (*serializers.RoleResp, *errors.RestErr)
	UpdateRole(roleID uint, req serializers.RoleReq) *errors.RestErr
	DeleteRole(id uint) *errors.RestErr
	ListRoles() ([]*serializers.RoleResp, *errors.RestErr)
	SetRolePermissions(req *serializers.RolePermissionsReq) *errors.RestErr
	GetRolePermissions(roleID int) ([]*domain.Permission, *errors.RestErr)
}
