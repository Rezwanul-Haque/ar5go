package impl

import (
	"clean/app/repository"
	"clean/infra/logger"
	"fmt"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type system struct {
	*gorm.DB
	*redis.Client
}

// NewSystemRepository will create an object that represent the System.Repository implementations
func NewSystemRepository(db *gorm.DB, redis *redis.Client) repository.ISystem {
	return &system{
		DB:     db,
		Client: redis,
	}
}

func (sys *system) DBCheck() (bool, error) {
	dB, _ := sys.DB.DB()
	if err := dB.Ping(); err != nil {
		return false, err
	}

	return true, nil
}

func (sys *system) CacheCheck() bool {
	client := sys.Client
	pong, err := client.Ping().Result()
	if err != nil {
		return false
	}

	logger.Info(fmt.Sprintf("%v from cache", pong))

	return true
}
