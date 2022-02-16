package repository

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infra/errors"
)

type IUsers interface {
	Save(user *domain.User) (*domain.User, *errors.RestErr)
	GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr)
	GetUserByID(userID uint) (*domain.User, *errors.RestErr)
	GetUserByEmail(email string) (*domain.User, error)
	Update(user *domain.User) *errors.RestErr
	UpdatePassword(userID uint, companyID uint, updateValue map[string]interface{}) *errors.RestErr
	GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr)
	GetUsersByCompanyIdAndRole(companyID, roleID uint,
		pagination *serializers.Pagination) ([]*domain.IntermediateUserResp, *errors.RestErr)
	SetLastLoginAt(user *domain.User) error
	HasRole(userID, roleID uint) bool
	ResetPassword(userID int, hashedPass []byte) error
	GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr)
	GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr)
}
