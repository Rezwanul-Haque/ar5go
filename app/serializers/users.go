package serializers

import (
	"clean/app/utils/consts"
	"clean/app/utils/methodutil"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type UserReq struct {
	UserName   string  `json:"user_name,omitempty"`
	FirstName  string  `json:"first_name,omitempty"`
	LastName   string  `json:"last_name,omitempty"`
	Email      string  `json:"email,omitempty"`
	Password   *string `json:"password,omitempty"`
	ProfilePic *string `json:"profile_pic,omitempty"`
	Phone      string  `json:"phone,omitempty"`
	CompanyID  uint    `json:"company_id"`
}

type LoggedInUser struct {
	ID          int      `json:"user_id"`
	AccessUuid  string   `json:"access_uuid"`
	RefreshUuid string   `json:"refresh_uuid"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type UserResp struct {
	ID          int        `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	ProfilePic  *string    `json:"profile_pic"`
	AppKey      string     `json:"app_key,omitempty"`
	RoleID      uint       `json:"role_id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login"`
}

type ResolveUserResp struct {
	ID          int        `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	ProfilePic  *string    `json:"profile_pic"`
	AppKey      string     `json:"app_key,omitempty"`
	RoleID      uint       `json:"role_id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login"`
}

type UserWithParamsResp struct {
	UserResp
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions,omitempty"`
}

type VerifyTokenResp struct {
	ID          int      `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email"`
	Phone       *string  `json:"phone"`
	ProfilePic  *string  `json:"profile_pic"`
	Permissions []string `json:"permissions"`
	Admin       bool     `json:"admin"`
}

type UserNameIsUnique struct {
	UserName string `json:"user_name"`
}

type EmailIsUnique struct {
	Email string `json:"email"`
}

func (lu LoggedInUser) IsAdmin() bool {
	return consts.RoleAdmin == lu.Role
}

type UserImagesResp struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
	Path    string `json:"path"`
}

func (lu LoggedInUser) HasPermission(perm string) bool {
	return methodutil.InArray(perm, lu.Permissions)
}

func (vr UserNameIsUnique) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.UserName, v.Required),
	)
}

func (vr EmailIsUnique) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.Email, v.Required),
	)
}
