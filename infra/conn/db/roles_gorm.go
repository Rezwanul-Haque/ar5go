package db

import (
	"ar5go/app/domain"
	"ar5go/app/utils/msgutil"
	"ar5go/infra/conn/db/models"
	"ar5go/infra/errors"
	"gorm.io/gorm"
)

func (dc DatabaseClient) CreateRole(role *domain.Role) (*domain.Role, *errors.RestErr) {
	res := dc.DB.Model(&models.Role{}).Where("name = ?", role.Name).FirstOrCreate(&role)

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("create role"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return role, nil
}

func (dc DatabaseClient) GetRole(roleID uint) (*domain.Role, *errors.RestErr) {
	var resp domain.Role

	res := dc.DB.Model(&models.Role{}).Where("id = ?", roleID).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("getting role by role id"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (dc DatabaseClient) UpdateRole(role *domain.Role) *errors.RestErr {
	res := dc.DB.Model(&models.Role{}).Where("id = ?", role.ID).Updates(&role)

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("update role"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) RemoveRole(id uint) *errors.RestErr {
	res := dc.DB.Where("id = ?", id).Delete(&models.Role{})

	if res.Error == gorm.ErrRecordNotFound {
		return errors.NewNotFoundError(msgutil.EntityNotFoundMsg("role"))
	}

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("remove role"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) ListRoles() ([]*domain.Role, *errors.RestErr) {
	var roles []*domain.Role

	res := dc.DB.Model(&models.Role{}).Find(&roles)

	if res.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewNotFoundError(msgutil.EntityNotFoundMsg("role"))
	}

	if res.Error != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("list roles"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return roles, nil
}

func (dc DatabaseClient) SetRolePermissions(rolePerms *domain.RolePermissions) *errors.RestErr {
	tx := dc.DB.Begin()

	if err := tx.Where("role_id = ?", rolePerms.RoleID).Delete(&domain.RolePermission{}).Error; err != nil {
		tx.Rollback()
		dc.lc.Error(msgutil.EntityGenericFailedMsg("deleting previous rule permissions"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	for _, permID := range rolePerms.Permissions {
		rp := &models.RolePermission{
			RoleID:       rolePerms.RoleID,
			PermissionID: permID,
		}

		if err := tx.Create(rp).Error; err != nil {
			tx.Rollback()
			dc.lc.Error(msgutil.EntityGenericFailedMsg("creating new rule permissions"), err)
			return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		}
	}

	if err := tx.Commit().Error; err != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("committing new rule permissions"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) GetRolePermissions(roleID int) ([]*domain.Permission, *errors.RestErr) {
	var res []*domain.Permission

	err := dc.DB.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&res).Error

	if err != nil {
		dc.lc.Error(msgutil.EntityGenericFailedMsg("getting role permissions by role id"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return res, nil
}
