package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/infra/conn/db"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
	"context"
)

type permissions struct {
	ctx context.Context
	lc  logger.LogClient
	DB  db.DatabaseClient
}

// NewPermissionsRepository will create an object that represent the Permission.Repository implementations
func NewPermissionsRepository(ctx context.Context, lc logger.LogClient, dbc db.DatabaseClient) repository.IPermissions {
	return &permissions{
		ctx: ctx,
		lc:  lc,
		DB:  dbc,
	}
}

func (r *permissions) CreatePermission(permission *domain.Permission) (*domain.Permission, *errors.RestErr) {
	return r.DB.CreatePermission(permission)
}

func (r *permissions) GetPermission(permissionID uint) (*domain.Permission, *errors.RestErr) {
	return r.DB.GetPermission(permissionID)
}

func (r *permissions) UpdatePermission(permission *domain.Permission) *errors.RestErr {
	return r.DB.UpdatePermission(permission)
}

func (r *permissions) RemovePermission(id uint) *errors.RestErr {
	return r.DB.RemovePermission(id)
}

func (r *permissions) ListPermissions() ([]*domain.Permission, *errors.RestErr) {
	return r.DB.ListPermissions()
}
