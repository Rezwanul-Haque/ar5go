package repository

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type Mails interface {
	SendUserCreatedEmail(mail *domain.UserCreateMail) *errors.RestErr
	SendCompanyCreatedEmail(mail *domain.CompanyCreatedMailReq) *errors.RestErr
	SendForgotPasswordEmail(mail *domain.ForgetPasswordMail) *errors.RestErr
}
