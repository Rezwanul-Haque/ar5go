package repository

import (
	"clean/app/domain"
	"clean/infrastructure/errors"
)

type IUsers interface {
	Save(user *domain.User) (*domain.User, *errors.RestErr)
	GetUser(userID uint) (*domain.User, error)
	GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr)
	GetUserByEmail(email string) (*domain.User, error)
	GetUsersByCompanyIdAndRole(companyID, roleID uint) ([]*domain.User, *errors.RestErr)
	SetLastLoginAt(user *domain.User) error
	HasRole(userID, roleID uint) bool
}
