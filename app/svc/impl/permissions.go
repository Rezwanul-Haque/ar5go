package impl

import (
	"boilerplate/app/domain"
	"boilerplate/app/repository"
	"boilerplate/app/serializers"
	"boilerplate/app/svc"
	"boilerplate/app/utils/methodutil"
	"boilerplate/app/utils/msgutil"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"
)

type permissions struct {
	prepo repository.IPermissions
}

func NewPermissionsService(prepo repository.IPermissions) svc.IPermissions {
	return &permissions{
		prepo: prepo,
	}
}

func (p *permissions) CreatePermission(req *serializers.PermissionReq) (*domain.Permission, *errors.RestErr) {
	permission := &domain.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	resp, err := p.prepo.Create(permission)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *permissions) GetPermission(id uint) (*serializers.PermissionResp, *errors.RestErr) {
	var resp serializers.PermissionResp
	rule, getErr := p.prepo.Get(id)
	if getErr != nil {
		return nil, getErr
	}

	err := methodutil.StructToStruct(rule, &resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("get permission"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return &resp, nil
}

func (p *permissions) UpdatePermission(permissionID uint, req serializers.PermissionReq) *errors.RestErr {
	permission := &domain.Permission{
		ID:          permissionID,
		Name:        req.Name,
		Description: req.Description,
	}

	upErr := p.prepo.Update(permission)
	if upErr != nil {
		return upErr
	}

	return nil
}

func (p *permissions) DeletePermission(id uint) *errors.RestErr {
	err := p.prepo.Remove(id)
	if err != nil {
		return err
	}

	return nil
}

func (p *permissions) ListPermission() ([]*serializers.PermissionResp, *errors.RestErr) {
	var resp []*serializers.PermissionResp

	permissions, lstErr := p.prepo.List()
	if lstErr != nil {
		return nil, lstErr
	}

	err := methodutil.StructToStruct(permissions, &resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("get all permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return resp, nil
}
