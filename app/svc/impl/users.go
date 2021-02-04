package svc

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/svc"
	"clean/infrastructure/errors"
)

type users struct {
	urepo repository.IUsers
}

func NewUsersService(urepo repository.IUsers) svc.IUsers {
	return &users{
		urepo: urepo,
	}
}

func (u *users) CreateAdminUser(user domain.User) (*domain.User, *errors.RestErr) {
	resp, saveErr := u.urepo.Save(&user)
	if saveErr != nil {
		return nil, saveErr
	}
	return resp, nil
}

func (u *users) CreateUser(user domain.User) (*domain.User, *errors.RestErr) {
	resp, saveErr := u.urepo.Save(&user)
	if saveErr != nil {
		return nil, saveErr
	}
	return resp, nil
}

func (u *users) GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr) {
	resp, getErr := u.urepo.GetUserByAppKey(appKey)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}

func (u *users) GetUserByCompanyIdAndRole(companyID, roleID uint) (domain.Users, *errors.RestErr) {
	resp, getErr := u.urepo.GetUsersByCompanyIdAndRole(companyID, roleID)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}
