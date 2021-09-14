package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/infra/config"
	"clean/infra/errors"
	"clean/infra/logger"
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type mails struct {
	mg *mailgun.MailgunImpl
}

// NewMailsRepository will create an object that represent the User.Repository implementations
func NewMailsRepository(mg *mailgun.MailgunImpl) repository.Mails {
	return &mails{
		mg: mg,
	}
}

func (r *mails) SendCompanyCreatedEmail(mail *domain.CompanyCreatedMailReq) *errors.RestErr {

	subject := GenerateEmailSubject(mail.To, config.Mail().Subject.UserCreated)
	ChangePass := fmt.Sprintf("Welcome to clean CRM! We're very excited to have you on board.\n Email: %v \n Password: %v \nThe password is auto generated, please change immediately. ", mail.To, mail.Password)
	m := r.mg.NewMessage(
		config.Mail().Sender,
		subject,
		ChangePass,
		mail.To,
		mail.Password,
	)

	m.SetTemplate(config.Mail().Mailgun.TemplateName.ForgotPassword) // this should be exists on https://app.mailgun.com/app/mail.shadowchef.co/templates
	m.AddTemplateVariable(config.Mail().PasswordResetUrl, ChangePass)
	addCommonTemplateVariable(m)

	if err := r.Send(m); err != nil {
		return err
	}
	return nil
}

func (r *mails) SendUserCreatedEmail(mail *domain.UserCreateMail) *errors.RestErr {

	subject := GenerateEmailSubject(mail.To, config.Mail().Subject.UserCreated)
	confirmationUrl := fmt.Sprintf("%v?id=%v&token=%v", config.Mail().UserConfirmationUrl, mail.UserID, mail.Token)

	m := r.mg.NewMessage(
		config.Mail().Sender,
		subject,
		confirmationUrl,
		mail.To,
		mail.Password,
	)

	m.SetTemplate(config.Mail().Mailgun.TemplateName.ForgotPassword) // this should be exists on https://app.mailgun.com/app/mail.shadowchef.co/templates
	m.AddTemplateVariable(config.Mail().PasswordResetUrl, confirmationUrl)
	addCommonTemplateVariable(m)

	if err := r.Send(m); err != nil {
		return err
	}

	return nil
}

func (r *mails) SendForgotPasswordEmail(mail *domain.ForgetPasswordMail) *errors.RestErr {

	subject := GenerateEmailSubject(mail.To, config.Mail().Subject.ForgotPassword)
	resetUrl := fmt.Sprintf("%v?id=%v&token=%v", config.Mail().PasswordResetUrl, mail.UserID, mail.Token)

	m := r.mg.NewMessage(
		config.Mail().Sender,
		subject,
		resetUrl,
		mail.To,
	)

	m.SetTemplate(config.Mail().Mailgun.TemplateName.ForgotPassword) // this should be exists on https://app.mailgun.com/app/mail.shadowchef.co/templates
	m.AddTemplateVariable(config.Mail().PasswordResetUrl, resetUrl)
	addCommonTemplateVariable(m)

	if err := r.Send(m); err != nil {
		return err
	}

	return nil
}

func addCommonTemplateVariable(m *mailgun.Message) {
	m.AddTemplateVariable("currentYear", time.Now().Year())
}

func (r *mails) Send(m *mailgun.Message) *errors.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := r.mg.Send(ctx, m)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func GenerateEmailSubject(entity, mailType string) string {
	subject := config.App().Name + " | " + entity + " | " + mailType
	logger.Info(fmt.Sprintf("Sending email '%s'...", subject))
	return subject
}
