package svc

import (
	"clean/app/svc"
	"clean/app/utils/methodsutil"
	"clean/infra/errors"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type redisCache struct {
	Redis *redis.Client
}

// NewRedisService will create an object that represent the Redis.Service implementations
func NewRedisService(redis *redis.Client) svc.ICache {
	return &redisCache{
		redis,
	}
}

func (rc *redisCache) Set(key string, value interface{}, ttl int) error {
	if methodsutil.IsEmpty(key) || methodsutil.IsEmpty(value) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rc.Redis.Set(key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
}

func (rc *redisCache) Get(key string) (string, error) {
	if methodsutil.IsEmpty(key) {
		return "", errors.ErrEmptyRedisKeyValue
	}

	return rc.Redis.Get(key).Result()
}

func (rc *redisCache) GetInt(key string) (int, error) {
	if methodsutil.IsEmpty(key) {
		return 0, errors.ErrEmptyRedisKeyValue
	}

	str, err := rc.Redis.Get(key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func (rc *redisCache) GetStruct(key string, outputStruct interface{}) error {
	if methodsutil.IsEmpty(key) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := rc.Redis.Get(key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(serializedValue), &outputStruct); err != nil {
		return err
	}

	return nil
}

func (rc *redisCache) Del(keys ...string) error {
	return rc.Redis.Del(keys...).Err()
}

func (rc *redisCache) DelPattern(pattern string) error {
	iter := rc.Redis.Scan(0, pattern, 0).Iterator()

	for iter.Next() {
		err := rc.Redis.Del(iter.Val()).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
