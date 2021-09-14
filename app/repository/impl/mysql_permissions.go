package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"

	"gorm.io/gorm"
)

type permissions struct {
	*gorm.DB
}

// NewMySqlPermissionsRepository will create an object that represent the Permission.Repository implementations
func NewMySqlPermissionsRepository(db *gorm.DB) repository.IPermissions {
	return &permissions{
		DB: db,
	}
}

func (r *permissions) Create(permission *domain.Permission) (*domain.Permission, *errors.RestErr) {
	res := r.DB.Model(&domain.Permission{}).Where("name = ?", permission.Name).FirstOrCreate(&permission)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("create permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return permission, nil
}

func (r *permissions) Get(permissionID uint) (*domain.Permission, *errors.RestErr) {
	var resp domain.Permission

	res := r.DB.Model(&domain.Permission{}).Where("id = ?", permissionID).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting permission by permission id"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (r *permissions) Update(permission *domain.Permission) *errors.RestErr {
	res := r.DB.Model(&domain.Permission{}).Where("id = ?", permission.ID).Updates(&permission)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("update permission"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *permissions) Remove(id uint) *errors.RestErr {
	res := r.DB.Where("id = ?", id).Delete(&domain.Permission{})

	if res.Error == gorm.ErrRecordNotFound {
		return errors.NewNotFoundError(msgutil.EntityNotFoundMsg("permission"))
	}

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("remove permission"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *permissions) List() ([]*domain.Permission, *errors.RestErr) {
	var permissions []*domain.Permission

	res := r.DB.Find(&permissions)

	if res.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewNotFoundError(msgutil.EntityNotFoundMsg("permissions"))
	}

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("list permissions"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return permissions, nil
}
