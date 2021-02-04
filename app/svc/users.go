package svc

import (
	"clean/app/domain"
	"clean/infrastructure/errors"
)

type IUsers interface {
	CreateAdminUser(domain.User) (*domain.User, *errors.RestErr)
	CreateUser(domain.User) (*domain.User, *errors.RestErr)
	GetUserByAppKey(apiKey string) (*domain.User, *errors.RestErr)
	GetUserByCompanyIdAndRole(companyID, roleID uint) (domain.Users, *errors.RestErr)
}
