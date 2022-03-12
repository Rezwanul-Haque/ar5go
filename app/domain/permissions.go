package domain

import "ar5go/infra/errors"

type IPermissions interface {
	CreatePermission(role *Permission) (*Permission, *errors.RestErr)
	GetPermission(roleID uint) (*Permission, *errors.RestErr)
	UpdatePermission(role *Permission) *errors.RestErr
	RemovePermission(id uint) *errors.RestErr
	ListPermissions() ([]*Permission, *errors.RestErr)
}

type Permission struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"unique" json:"name"`
	Description string `json:"description"`
}
