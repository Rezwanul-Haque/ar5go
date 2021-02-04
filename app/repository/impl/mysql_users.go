package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/infrastructure/errors"
	"clean/infrastructure/logger"
	"gorm.io/gorm"
	"time"
)

type users struct {
	*gorm.DB
}

// NewMySqlUsersRepository will create an object that represent the User.Repository implementations
func NewMySqlUsersRepository(db *gorm.DB) repository.IUsers {
	return &users{
		DB: db,
	}
}

func (r *users) Save(user *domain.User) (*domain.User, *errors.RestErr) {
	res := r.DB.Model(&domain.User{}).Create(&user)

	if res.Error != nil {
		logger.Error("error occurred when create user", res.Error)
		return nil, errors.NewInternalServerError("db error")
	}

	return user, nil
}

func (r *users) GetUser(userID uint) (*domain.User, error) {
	var resp domain.User

	res := r.DB.Model(&domain.User{}).Where("id = ?", userID).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.ErrRecordNotFound
	}

	if res.Error != nil {
		logger.Error("error occurred when getting user by app key", res.Error)
		return nil, errors.NewError("db error")
	}

	return &resp, nil
}

func (r *users) GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr) {
	var resp domain.User

	res := r.DB.Model(&domain.User{}).Where("app_key = ?", appKey).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("no user found")
	}

	if res.Error != nil {
		logger.Error("error occurred when getting user by app key", res.Error)
		return nil, errors.NewInternalServerError("db error")
	}

	return &resp, nil
}

func (r *users) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	res := r.DB.Model(&domain.User{}).Where("email = ?", email).Find(&user)
	if res.RowsAffected == 0 {
		logger.Error("no user found by this email", res.Error)
		return nil, errors.ErrRecordNotFound
	}
	if res.Error != nil {
		logger.Error("error occurred when trying to get user by email", res.Error)
		return nil, errors.NewError("db error")
	}

	return user, nil
}

func (r *users) GetUsersByCompanyIdAndRole(companyID, roleID uint) ([]*domain.User, *errors.RestErr) {
	var resp []*domain.User

	res := r.DB.Model(&domain.User{}).Where("company_id = ? AND role_id = ?", companyID, roleID).Find(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("no users found")
	}

	if res.Error != nil {
		logger.Error("error occurred when getting users by company_id and role_id", res.Error)
		return nil, errors.NewInternalServerError("db error")
	}

	return resp, nil
}

func (r *users) SetLastLoginAt(user *domain.User) error {
	lastLoginAt := time.Now().UTC()

	err := r.DB.Model(&user).Update("last_login_at", lastLoginAt).Error

	if err != nil {
		logger.Error(err.Error(), err)
		return err
	}

	return nil
}

func (r *users) HasRole(userID, roleID uint) bool {
	var count int64
	count = 0

	r.DB.Model(&domain.User{}).
		Where("id = ? AND role_id = ?", userID, roleID).
		Count(&count)

	return count > 0
}
