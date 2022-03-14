package db

import (
	"ar5go/app/domain"
	"ar5go/app/utils/msgutil"
	"ar5go/infra/conn/db/models"
	"ar5go/infra/errors"
	"gorm.io/gorm"
)

func (dc DatabaseClient) CreatePermission(permission *domain.Permission) (*domain.Permission, *errors.RestErr) {
	res := dc.DB.Model(&models.Permission{}).Where("name = ?", permission.Name).FirstOrCreate(&permission)

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("create permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return permission, nil
}

func (dc DatabaseClient) GetPermission(permissionID uint) (*domain.Permission, *errors.RestErr) {
	var resp domain.Permission

	res := dc.DB.Model(&models.Permission{}).Where("id = ?", permissionID).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("getting permission by permission id"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (dc DatabaseClient) UpdatePermission(permission *domain.Permission) *errors.RestErr {
	res := dc.DB.Model(&models.Permission{}).Where("id = ?", permission.ID).Updates(&permission)

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("update permission"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) RemovePermission(id uint) *errors.RestErr {
	res := dc.DB.Where("id = ?", id).Delete(&models.Permission{})

	if res.Error == gorm.ErrRecordNotFound {
		return errors.NewNotFoundError(msgutil.EntityNotFoundMsg("permission"))
	}

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("remove permission"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) ListPermissions() ([]*domain.Permission, *errors.RestErr) {
	var permissions []*domain.Permission

	res := dc.DB.Model(&models.Permission{}).Find(&permissions)

	if res.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewNotFoundError(msgutil.EntityNotFoundMsg("permissions"))
	}

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("list permissions"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return permissions, nil
}
