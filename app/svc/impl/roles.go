package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/app/serializers"
	"ar5go/app/svc"
	"ar5go/app/utils/methodsutil"
	"ar5go/app/utils/msgutil"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
)

type roles struct {
	rrepo repository.IRoles
}

func NewRolesService(rrepo repository.IRoles) svc.IRoles {
	return &roles{
		rrepo: rrepo,
	}
}

func (r *roles) CreateRole(req *serializers.RoleReq) (*domain.Role, *errors.RestErr) {
	role := &domain.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
	}

	resp, err := r.rrepo.CreateRole(role)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *roles) GetRole(id uint) (*serializers.RoleResp, *errors.RestErr) {
	var resp *serializers.RoleResp
	rule, getErr := r.rrepo.GetRole(id)
	if getErr != nil {
		return nil, getErr
	}

	err := methodsutil.StructToStruct(rule, &resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("get role"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return resp, nil
}

func (r *roles) UpdateRole(roleID uint, req serializers.RoleReq) *errors.RestErr {
	role := &domain.Role{
		ID:          roleID,
		Name:        req.Name,
		DisplayName: req.DisplayName,
	}

	upErr := r.rrepo.UpdateRole(role)
	if upErr != nil {
		return upErr
	}

	return nil
}

func (r *roles) DeleteRole(id uint) *errors.RestErr {
	err := r.rrepo.RemoveRole(id)
	if err != nil {
		return err
	}

	return nil
}

func (r *roles) ListRoles() ([]*serializers.RoleResp, *errors.RestErr) {
	var resp []*serializers.RoleResp

	rules, lstErr := r.rrepo.ListRoles()
	if lstErr != nil {
		return nil, lstErr
	}

	err := methodsutil.StructToStruct(rules, &resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("get role"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return resp, nil
}

func (r *roles) SetRolePermissions(req *serializers.RolePermissionsReq) *errors.RestErr {
	var rolePerms *domain.RolePermissions

	err := methodsutil.StructToStruct(req, &rolePerms)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set role & permissions"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	setErr := r.rrepo.SetRolePermissions(rolePerms)
	if setErr != nil {
		return setErr
	}
	return nil
}

func (r *roles) GetRolePermissions(roleID int) ([]*domain.Permission, *errors.RestErr) {
	resp, getErr := r.rrepo.GetRolePermissions(roleID)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}
