package svc

import (
	"boilerplate/app/serializers"
	"boilerplate/infra/errors"
)

type IMails interface {
	SendCompanyCreatedEmail(req serializers.CompanyCreatedMailReq) *errors.RestErr
	SendUserCreatedEmail(req serializers.UserCreatedMailReq) *errors.RestErr
	SendForgotPasswordEmail(req serializers.ForgetPasswordMailReq) *errors.RestErr
}
