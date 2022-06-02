package svc

import (
	"ar5go/app/domain"
	"ar5go/app/serializers"
	"ar5go/infra/errors"
)

type IUsers interface {
	CreateAdminUser(domain.User) (*domain.User, *errors.RestErr)
	CreateUser(domain.User) (*domain.User, *errors.RestErr)
	GetUserById(uid uint) (*domain.User, *errors.RestErr)
	GetUserByEmail(useremail string) (*domain.User, error)
	GetUserByAppKey(apiKey string) (*domain.User, *errors.RestErr)
	UpdateUser(userID uint, req serializers.UserReq) *errors.RestErr
	GetUserByCompanyIdAndRole(companyID, roleID uint, filters *serializers.ListFilters) (*serializers.ListFilters, *errors.RestErr)
	ChangePassword(id int, data *serializers.ChangePasswordReq) error
	ForgotPassword(email string) error
	VerifyResetPassword(req *serializers.VerifyResetPasswordReq) error
	ResetPassword(req *serializers.ResetPasswordReq) error
}
