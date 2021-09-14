package svc

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infra/errors"
)

type IPermissions interface {
	CreatePermission(req *serializers.PermissionReq) (*domain.Permission, *errors.RestErr)
	GetPermission(id uint) (*serializers.PermissionResp, *errors.RestErr)
	UpdatePermission(permissionID uint, req serializers.PermissionReq) *errors.RestErr
	DeletePermission(id uint) *errors.RestErr
	ListPermission() ([]*serializers.PermissionResp, *errors.RestErr)
}
