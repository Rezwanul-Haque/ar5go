package repository

import (
	"boilerplate/app/domain"
	"boilerplate/infra/errors"
)

type IUsers interface {
	Save(user *domain.User) (*domain.User, *errors.RestErr)
	GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr)
	GetUserByID(userID uint) (*domain.User, *errors.RestErr)
	UserNameIsUnique(username string) error
	EmailIsUnique(email string) error
	GetUserByEmail(email string) (*domain.User, error)
	Update(user *domain.User) *errors.RestErr
	UpdateUserActivation(id uint, activate bool) *errors.RestErr
	UpdatePassword(userID uint, updateValue map[string]interface{}) *errors.RestErr
	SetLastLoginAt(user *domain.User) error
	HasRole(userID, roleID uint) bool
	ResetPassword(userID int, hashedPass []byte) error
	GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr)
	GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr)
}
