package cache

import (
	"ar5go/infra/config"
	"ar5go/infra/logger"
	"context"
	"github.com/go-redis/redis/v8"
)

type CacheClient struct {
	Redis *redis.Client
}

func connectRedis() {
	conf := config.Cache().Redis

	logger.Info("connecting to Redis at " + conf.Host + ":" + conf.Port + "...")

	c := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Pass,
		DB:       conf.Db,
	})

	client.Redis = c

	if _, err := client.Redis.Ping(context.Background()).Result(); err != nil {
		logger.Error("failed to connect Redis: ", err)
		panic(err)
	}

	logger.Info("Redis connection successful...")
}
