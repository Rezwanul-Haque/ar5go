package impl

import (
	"boilerplate/app/domain"
	"boilerplate/app/repository"
	"boilerplate/app/serializers"
	"boilerplate/app/svc"
	"boilerplate/infra/errors"
)

type mails struct {
	mrepo repository.Mails
}

func NewMailsService(mrepo repository.Mails) svc.IMails {
	return &mails{
		mrepo: mrepo,
	}
}

func (m *mails) SendCompanyCreatedEmail(req serializers.CompanyCreatedMailReq) *errors.RestErr {
	mail := &domain.CompanyCreatedMailReq{
		To:        req.To,
		CompanyID: req.CompanyID,
		Password:  req.Password,
		Token:     req.Token,
	}

	if err := m.mrepo.SendCompanyCreatedEmail(mail); err != nil {
		return err
	}
	return nil
}

func (m *mails) SendUserCreatedEmail(req serializers.UserCreatedMailReq) *errors.RestErr {
	mail := &domain.UserCreateMail{
		To:       req.To,
		UserID:   req.UserID,
		Password: req.Password,
		Token:    req.Token,
	}

	if err := m.mrepo.SendUserCreatedEmail(mail); err != nil {
		return err
	}
	return nil
}

func (m *mails) SendForgotPasswordEmail(req serializers.ForgetPasswordMailReq) *errors.RestErr {
	mail := &domain.ForgetPasswordMail{
		To:     req.To,
		UserID: req.UserID,
		Token:  req.Token,
	}

	if err := m.mrepo.SendForgotPasswordEmail(mail); err != nil {
		return err
	}
	return nil
}
