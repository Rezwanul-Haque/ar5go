package domain

import (
	"boilerplate/app/utils/methodutil"
	"boilerplate/infra/errors"
	"mime/multipart"
	"strconv"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID          uint       `json:"id"`
	UserName    string     `json:"user_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `gorm:"unique" json:"email"`
	Password    *string    `json:"password,omitempty"`
	RoleID      uint       `json:"role_id"`
	ProfilePic  *string    `json:"profile_pic"`
	LastLoginAt *time.Time `json:"last_login_at"`
	FirstLogin  bool       `json:"first_login" gorm:"column:first_login;default:true"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

type Users []*User

type IntermediateUserWithPermissions struct {
	User
	RoleName    string `json:"role_name"`
	Permissions string `json:"permissions"`
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

type UserImage struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Path   string `json:"path"`
}

type UserImageStorage struct {
	Image            *multipart.FileHeader
	HashedImagesName string
	Path             string
	BasePath         string
}

func (uimg *UserImageStorage) GenerateBasePath(fileType string, userID int) {
	if fileType == "img" {
		uimg.BasePath = "images/users/" + strconv.Itoa(userID) + "/"
	}
}

func (uimg *UserImageStorage) Validate(Type string) *errors.RestErr {
	if err := methodutil.ValidateImageFileType(uimg.Image, Type); err != nil {
		return err
	}

	return nil
}

func (vr User) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.UserName, v.Required),
		v.Field(&vr.Email, v.Required, is.Email),
	)
}
