package repository

import (
	"boilerplate/app/domain"
	"boilerplate/infra/errors"
)

type Mails interface {
	SendUserCreatedEmail(mail *domain.UserCreateMail) *errors.RestErr
	SendCompanyCreatedEmail(mail *domain.CompanyCreatedMailReq) *errors.RestErr
	SendForgotPasswordEmail(mail *domain.ForgetPasswordMail) *errors.RestErr
}
