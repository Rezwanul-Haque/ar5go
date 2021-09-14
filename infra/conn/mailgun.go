package conn

import (
	"boilerplate/infra/config"
	"boilerplate/infra/logger"

	"github.com/mailgun/mailgun-go/v4"
)

var mg *mailgun.MailgunImpl

func ConnectMailGun() {
	conf := config.Mail()
	mg = mailgun.NewMailgun(conf.Mailgun.Domain, conf.Mailgun.ApiKey)

	logger.Info("mailgun client created...")
}

func MailGun() *mailgun.MailgunImpl {
	return mg
}
