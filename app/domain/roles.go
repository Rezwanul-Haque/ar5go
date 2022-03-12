package domain

import "ar5go/infra/errors"

type IRoles interface {
	CreateRole(role *Role) (*Role, *errors.RestErr)
	GetRole(roleID uint) (*Role, *errors.RestErr)
	UpdateRole(role *Role) *errors.RestErr
	RemoveRole(id uint) *errors.RestErr
	ListRoles() ([]*Role, *errors.RestErr)
	SetRolePermissions(rolePerms *RolePermissions) *errors.RestErr
	GetRolePermissions(roleID int) ([]*Permission, *errors.RestErr)
}

type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"unique" json:"name"`
	DisplayName string `json:"display_name"`
}
