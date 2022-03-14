package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/app/serializers"
	"ar5go/infra/conn/db"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
	"context"
	"time"
)

type users struct {
	ctx context.Context
	lc  logger.LogClient
	DB  db.DatabaseClient
}

// NewUsersRepository will create an object that represent the User.Repository implementations
func NewUsersRepository(ctx context.Context, lc logger.LogClient, dbc db.DatabaseClient) repository.IUsers {
	return &users{
		ctx: ctx,
		lc:  lc,
		DB:  dbc,
	}
}

func (r *users) SaveUser(user *domain.User) (*domain.User, *errors.RestErr) {
	return r.DB.SaveUser(user)
}

func (r *users) GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	return r.DB.GetUser(userID, withPermission)
}

func (r *users) GetUserByID(userID uint) (*domain.User, *errors.RestErr) {
	return r.DB.GetUserByID(userID)
}

func (r *users) UpdateUser(user *domain.User) *errors.RestErr {
	return r.DB.UpdateUser(user)
}

func (r *users) UpdatePassword(userID uint, companyID uint, updateValues map[string]interface{}) *errors.RestErr {
	return r.DB.UpdatePassword(userID, companyID, updateValues)
}

func (r *users) GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr) {
	return r.DB.GetUserByAppKey(appKey)
}

func (r *users) GetUserByEmail(email string) (*domain.User, error) {
	return r.DB.GetUserByEmail(email)
}

func (r *users) GetUsersByCompanyIdAndRole(companyID, roleID uint,
	filters *serializers.ListFilters) ([]*domain.IntermediateUserResp, *errors.RestErr) {
	return r.DB.GetUsersByCompanyIdAndRole(companyID, roleID, filters)
}

func (r *users) SetLastLoginAt(user *domain.User) error {
	*user.LastLoginAt = time.Now().UTC()

	return r.DB.SetLastLoginAt(user)
}

func (r *users) HasRole(userID, roleID uint) bool {
	return r.DB.HasRole(userID, roleID)
}

func (r *users) ResetPassword(userID int, hashedPass []byte) error {
	return r.DB.ResetPassword(userID, hashedPass)
}

func (r *users) GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr) {
	return r.DB.GetTokenUser(id)
}

func (r *users) GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	return r.DB.GetUserWithPermissions(userID, withPermission)
}
