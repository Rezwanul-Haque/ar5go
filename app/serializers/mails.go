package serializers

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type MailReq struct {
	To       string `json:"to" valid:"required"`
	Subject  string `json:"subject"`
	Password string `json:"password" valid:"required"`
}

type CompanyCreatedMailReq struct {
	To        string `json:"to" valid:"required"`
	CompanyID uint   `json:"user_id" valid:"required"`
	Password  string `json:"password" valid:"required"`
	Token     string `json:"token"`
}

type UserCreatedMailReq struct {
	To       string `json:"to" valid:"required"`
	UserID   uint   `json:"user_id" valid:"required"`
	Password string `json:"password" valid:"required"`
	Token    string `json:"token"`
}

type ForgetPasswordMailReq struct {
	To     string `json:"to" valid:"required"`
	UserID uint   `json:"user_id" valid:"required"`
	Token  string `json:"token"`
}

type GenericMailReq struct {
	To      string `json:"to" valid:"required"`
	Subject string `json:"subject" valid:"required"`
	Body    string `json:"body" valid:"-"`
}

func (m *MailReq) Validate() error {
	return v.ValidateStruct(m,
		v.Field(&m.To, is.EmailFormat),
	)
}

func (m *ForgetPasswordMailReq) Validate() error {
	return v.ValidateStruct(m,
		v.Field(&m.To, is.EmailFormat),
	)
}

func (m *GenericMailReq) Validate() error {
	return v.ValidateStruct(m,
		v.Field(&m.To, is.EmailFormat),
	)
}
