package repository

import (
	"boilerplate/app/domain"
	"boilerplate/infra/errors"
)

type IPermissions interface {
	Create(role *domain.Permission) (*domain.Permission, *errors.RestErr)
	Get(roleID uint) (*domain.Permission, *errors.RestErr)
	Update(role *domain.Permission) *errors.RestErr
	Remove(id uint) *errors.RestErr
	List() ([]*domain.Permission, *errors.RestErr)
}
