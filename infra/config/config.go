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
	Name string
	Port string
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

type RedisConfig struct {
	Host              string
	Port              string
	Pass              string
	Db                int
	AccessUuidPrefix  string
	RefreshUuidPrefix string
	UserPrefix        string
	TokenPrefix       string
	Ttl               int // seconds
}

type Config struct {
	App   *AppConfig
	Db    *DbConfig
	Jwt   *JwtConfig
	Redis *RedisConfig
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

func Redis() *RedisConfig {
	return config.Redis
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
		Name: "CLEAN",
		Port: "8080",
	}

	config.Db = &DbConfig{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "root",
		Pass:            "12345678",
		Schema:          "clean_db",
		MaxIdleConn:     1,
		MaxOpenConn:     2,
		MaxConnLifetime: 30,
		Debug:           true,
	}

	config.Jwt = &JwtConfig{
		AccessTokenSecret:  "accesstokensecret",
		RefreshTokenSecret: "refreshtokensecret",
		AccessTokenExpiry:  300,
		RefreshTokenExpiry: 10080,
		ContextKey:         "user",
	}

	config.Redis = &RedisConfig{
		Host:              "127.0.0.1",
		Port:              "6379",
		Pass:              "",
		Db:                0,
		AccessUuidPrefix:  "access-uuid_",
		RefreshUuidPrefix: "refresh-uuid_",
		UserPrefix:        "user_",
		TokenPrefix:       "token_",
		Ttl:               3600,
	}
}
