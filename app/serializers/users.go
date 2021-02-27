package serializers

import (
	"clean/app/domain"
	"clean/app/utils/consts"
	"clean/app/utils/methodsutil"
	"time"
)

type ResolveResp struct {
	CompanyName  string       `json:"company_name"`
	CompanyID    uint         `json:"company_id"`
	Subordinates domain.Users `json:"subordinates"`
}

type UserReq struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
	Phone     string  `json:"phone"`
}

type LoggedInUser struct {
	ID          int      `json:"user_id"`
	AccessUuid  string   `json:"access_uuid"`
	RefreshUuid string   `json:"refresh_uuid"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type UserResp struct {
	ID           int        `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	Phone        *string    `json:"phone"`
	ProfilePic   *string    `json:"profile_pic"`
	BusinessID   *int       `json:"business_id"`
	AppKey       string     `json:"app_key,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	FirstLogin   bool       `json:"first_login"`
}

type VerifyTokenResp struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	Phone      *string `json:"phone"`
	ProfilePic *string `json:"profile_pic"`
	BusinessID *int    `json:"business_id"`
	Admin      bool    `json:"admin"`
}

func (lu LoggedInUser) IsAdmin() bool {
	return methodsutil.InArray(consts.RoleAdmin, lu.Roles)
}

func (lu LoggedInUser) IsSales() bool {
	return methodsutil.InArray(consts.RoleSales, lu.Roles)
}
