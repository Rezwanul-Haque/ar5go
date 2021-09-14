package impl

import (
	"boilerplate/app/domain"
	"boilerplate/app/repository"
	"boilerplate/app/utils/msgutil"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"

	"gorm.io/gorm"
)

type roles struct {
	*gorm.DB
}

// NewMySqlRolesRepository will create an object that represent the Roles.Repository implementations
func NewMySqlRolesRepository(db *gorm.DB) repository.IRoles {
	return &roles{
		DB: db,
	}
}

func (r *roles) Create(role *domain.Role) (*domain.Role, *errors.RestErr) {
	res := r.DB.Model(&domain.Role{}).Where("name = ?", role.Name).FirstOrCreate(&role)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("create role"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return role, nil
}

func (r *roles) Get(roleID uint) (*domain.Role, *errors.RestErr) {
	var resp domain.Role

	res := r.DB.Model(&domain.Role{}).Where("id = ?", roleID).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting role by role id"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (r *roles) Update(role *domain.Role) *errors.RestErr {
	res := r.DB.Model(&domain.Role{}).Where("id = ?", role.ID).Updates(&role)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("update role"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *roles) Remove(id uint) *errors.RestErr {
	res := r.DB.Where("id = ?", id).Delete(&domain.Role{})

	if res.Error == gorm.ErrRecordNotFound {
		return errors.NewNotFoundError(msgutil.EntityNotFoundMsg("role"))
	}

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("remove role"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *roles) List() ([]*domain.Role, *errors.RestErr) {
	var rules []*domain.Role

	res := r.DB.Find(&rules)

	if res.Error == gorm.ErrRecordNotFound {
		return nil, errors.NewNotFoundError(msgutil.EntityNotFoundMsg("role"))
	}

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("list roles"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return rules, nil
}

func (r *roles) SetRolePermissions(rolePerms *domain.RolePermissions) *errors.RestErr {
	tx := r.DB.Begin()

	if err := tx.Where("role_id = ?", rolePerms.RoleID).Delete(&domain.RolePermission{}).Error; err != nil {
		tx.Rollback()
		logger.Error(msgutil.EntityGenericFailedMsg("deleting previous rule permissions"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	for _, permID := range rolePerms.Permissions {
		rp := &domain.RolePermission{
			RoleID:       rolePerms.RoleID,
			PermissionID: permID,
		}

		if err := tx.Create(rp).Error; err != nil {
			tx.Rollback()
			logger.Error(msgutil.EntityGenericFailedMsg("creating new rule permissions"), err)
			return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("committing new rule permissions"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *roles) GetRolePermissions(roleID int) ([]*domain.Permission, *errors.RestErr) {
	var res []*domain.Permission

	err := r.DB.
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&res).Error

	if err != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting role permissions by role id"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return res, nil
}
