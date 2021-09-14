package users

import (
	"clean/app/domain"
	"clean/infra/config"
	"clean/infra/errors"

	"github.com/stretchr/testify/mock"
)

func init() {
	setup()
}

func setup() {
	config.LoadConfig()
}

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) Save(user *domain.User) (*domain.User, *errors.RestErr) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*domain.User), nil
}
func (mock *MockRepository) GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*domain.UserWithPerms), nil
}
func (mock *MockRepository) GetUserByID(userID uint) (*domain.User, *errors.RestErr) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*domain.User), nil
}
func (mock *MockRepository) UserNameIsUnique(username string) error {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) PhoneNoIsUnique(phone string) error {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) EmailIsUnique(email string) error {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return result.(*domain.User), nil
}
func (mock *MockRepository) Update(user *domain.User) *errors.RestErr {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) UpdateUserActivation(id uint, activate bool) *errors.RestErr {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) UpdatePassword(userID uint, updateValue map[string]interface{}) *errors.RestErr {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*domain.User), nil
}

func (mock *MockRepository) SetLastLoginAt(user *domain.User) error {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) HasRole(userID, roleID uint) bool {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return args.Bool(0)
}
func (mock *MockRepository) ResetPassword(userID int, hashedPass []byte) error {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return nil
}
func (mock *MockRepository) GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return result.(*domain.UserWithPerms), nil
}
func (mock *MockRepository) GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr) {
	args := mock.Called()
	result := args.Get(0)
	_ = result
	return result.(*domain.VerifyTokenResp), nil
}
