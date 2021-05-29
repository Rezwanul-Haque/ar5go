package domain

import (
	"time"
)

type User struct {
	ID          uint       `json:"id"`
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
