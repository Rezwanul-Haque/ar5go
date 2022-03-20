package domain

import (
	"ar5go/app/serializers"
	"ar5go/infra/errors"
	"time"
)

type IUsers interface {
	SaveUser(user *User) (*User, *errors.RestErr)
	GetUser(userID uint, withPermission bool) (*UserWithPerms, *errors.RestErr)
	GetUserByID(userID uint) (*User, *errors.RestErr)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) *errors.RestErr
	UpdatePassword(userID uint, companyID uint, updateValue map[string]interface{}) *errors.RestErr
	GetUserByAppKey(appKey string) (*User, *errors.RestErr)
	GetUsersByCompanyIdAndRole(companyID, roleID uint,
		filters *serializers.ListFilters) ([]*IntermediateUserResp, *errors.RestErr)
	SetLastLoginAt(user *User) error
	HasRole(userID, roleID uint) bool
	ResetPassword(userID int, hashedPass []byte) error
	GetUserWithPermissions(userID uint, withPermission bool) (*UserWithPerms, *errors.RestErr)
	GetTokenUser(id uint) (*VerifyTokenResp, *errors.RestErr)
}

type User struct {
	ID          uint       `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Password    *string    `json:"password,omitempty"`
	Phone       string     `json:"phone"`
	CompanyID   uint       `json:"company_id"`
	AppKey      string     `json:"app_key,omitempty"`
	RoleID      uint       `json:"role_id"`
	ProfilePic  *string    `json:"profile_pic"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login" gorm:"column:first_login;default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type Users []*User

type IntermediateUserWithPermissions struct {
	User
	RoleName    string `json:"role_name"`
	Permissions string `json:"permissions"`
}

type IntermediateUserResp struct {
	CompanyID   uint       `json:"company_id"`
	CompanyName string     `json:"company_name"`
	ID          int        `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	ProfilePic  *string    `json:"profile_pic"`
	RoleID      uint       `json:"role_id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login"`
}

type UserWithPerms struct {
	User
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions,omitempty"`
}

type BaseVerifyTokenResp struct {
	ID           int     `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Phone        *string `json:"phone"`
	ProfilePic   *string `json:"profile_pic"`
	BusinessID   *int    `json:"business_id"`
	BusinessName string  `json:"business_name"`
	CompanyID    *int    `json:"company_id"`
	CompanyName  string  `json:"company_name"`
	Admin        bool    `json:"admin"`
	SuperAdmin   bool    `json:"super_admin"`
}

type TempVerifyTokenResp struct {
	BaseVerifyTokenResp
	Permissions string `json:"permissions"`
}

type VerifyTokenResp struct {
	BaseVerifyTokenResp
	Permissions []string `json:"permissions"`
}
