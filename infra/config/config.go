package config

import (
	"clean/infra/logger"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type AppConfig struct {
	Name                   string
	Port                   string
	Page                   int64
	Limit                  int64
	Sort                   string
	UserConfirmationPrefix string
	UserCreateAuthSkip     bool
	MockPasswordEnabled    bool
	MockPassword           string
	SendEmail              bool
}

type DbConfig struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Schema          string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
	Debug           bool
}

type JwtConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	ContextKey         string
}

type TemplateNameConfig struct {
	UserCreate     string
	ForgotPassword string
}

type MailgunConfig struct {
	ApiKey       string
	Domain       string
	TemplateName TemplateNameConfig
}

type MailgunSubject struct {
	UserCreated    string
	ForgotPassword string
	PasswordReset  string
}

type MailConfig struct {
	Sender              string
	PasswordResetUrl    string
	UserConfirmationUrl string
	Subject             MailgunSubject
	Mailgun             MailgunConfig
}

type Config struct {
	App  *AppConfig
	Db   *DbConfig
	Jwt  *JwtConfig
	Mail *MailConfig
}

var config Config

func App() *AppConfig {
	return config.App
}

func Db() *DbConfig {
	return config.Db
}

func Jwt() *JwtConfig {
	return config.Jwt
}

func Mail() *MailConfig {
	return config.Mail
}

func LoadConfig() {
	setDefaultConfig()

	_ = viper.BindEnv("CONSUL_URL")
	_ = viper.BindEnv("CONSUL_PATH")

	consulURL := viper.GetString("CONSUL_URL")
	consulPath := viper.GetString("CONSUL_PATH")

	if consulURL != "" && consulPath != "" {
		_ = viper.AddRemoteProvider("consul", consulURL, consulPath)

		viper.SetConfigType("json")
		err := viper.ReadRemoteConfig()

		if err != nil {
			log.Println(fmt.Sprintf("%s named \"%s\"", err.Error(), consulPath))
		}

		config = Config{}

		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}

		if r, err := json.MarshalIndent(&config, "", "  "); err == nil {
			fmt.Println(string(r))
		}
	} else {
		logger.Info("CONSUL_URL or CONSUL_PATH missing! Serving with default config...")
	}
}

func setDefaultConfig() {
	config.App = &AppConfig{
		Name:                   "clean",
		Port:                   "8080",
		Page:                   1,
		Limit:                  10,
		Sort:                   "created_at desc",
		UserConfirmationPrefix: "user-activation_",
		UserCreateAuthSkip:     false,
		MockPasswordEnabled:    true,
		MockPassword:           "12345678",
		SendEmail:              false,
	}

	config.Db = &DbConfig{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "root",
		Pass:            "12345678",
		Schema:          "clean",
		MaxIdleConn:     1,
		MaxOpenConn:     2,
		MaxConnLifetime: 30,
		Debug:           true,
	}

	config.Jwt = &JwtConfig{
		AccessTokenSecret:  "accesstokensecret",
		RefreshTokenSecret: "refreshtokensecret",
		AccessTokenExpiry:  3000,
		RefreshTokenExpiry: 10080,
		ContextKey:         "user",
	}

	config.Mail = &MailConfig{
		Sender:           "clean Admin <admin@clean.xyz>",
		PasswordResetUrl: "http://localhost:3000/password/reset",
		Subject: MailgunSubject{
			UserCreated:    "User Created",
			ForgotPassword: "Forgot Password",
			PasswordReset:  "Password Reset",
		},
		Mailgun: MailgunConfig{
			ApiKey: "mail-api-key",
			Domain: "mail.domain.xyz",
			TemplateName: TemplateNameConfig{
				UserCreate:     "Welcome",
				ForgotPassword: "password_reset",
			},
		},
	}
}
