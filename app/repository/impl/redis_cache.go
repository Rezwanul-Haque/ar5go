package impl

import (
	"clean/app/repository"
	"clean/app/utils/methodsutil"
	"clean/infra/errors"
	"encoding/json"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type redisCache struct {
	Redis *redis.Client
}

// NewRedisRepository will create an object that represent the Redis.Repository implementations
func NewRedisRepository(redis *redis.Client) repository.ICache {
	return &redisCache{
		redis,
	}
}

func (rRepo *redisCache) Set(key string, value interface{}, ttl int) error {
	if methodsutil.IsEmpty(key) || methodsutil.IsEmpty(value) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rRepo.Redis.Set(key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
}

func (rRepo *redisCache) Get(key string) (string, error) {
	if methodsutil.IsEmpty(key) {
		return "", errors.ErrEmptyRedisKeyValue
	}

	return rRepo.Redis.Get(key).Result()
}

func (rRepo *redisCache) GetInt(key string) (int, error) {
	if methodsutil.IsEmpty(key) {
		return 0, errors.ErrEmptyRedisKeyValue
	}

	str, err := rRepo.Redis.Get(key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func (rRepo *redisCache) GetStruct(key string, outputStruct interface{}) error {
	if methodsutil.IsEmpty(key) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := rRepo.Redis.Get(key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(serializedValue), &outputStruct); err != nil {
		return err
	}

	return nil
}

func (rRepo *redisCache) Del(keys ...string) error {
	return rRepo.Redis.Del(keys...).Err()
}

func (rRepo *redisCache) DelPattern(pattern string) error {
	iter := rRepo.Redis.Scan(0, pattern, 0).Iterator()

	for iter.Next() {
		err := rRepo.Redis.Del(iter.Val()).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
