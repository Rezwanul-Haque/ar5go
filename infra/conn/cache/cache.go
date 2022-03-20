package cache

import (
	"ar5go/app/utils/methodsutil"
	"ar5go/infra/errors"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

func (cc CacheClient) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	if methodsutil.IsEmpty(key) || methodsutil.IsEmpty(value) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cc.Redis.Set(ctx, key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
}

func (cc CacheClient) Get(ctx context.Context, key string) (string, error) {
	if methodsutil.IsEmpty(key) {
		return "", errors.ErrEmptyRedisKeyValue
	}

	return cc.Redis.Get(ctx, key).Result()
}

func (cc CacheClient) GetInt(ctx context.Context, key string) (int, error) {
	if methodsutil.IsEmpty(key) {
		return 0, errors.ErrEmptyRedisKeyValue
	}

	str, err := cc.Redis.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func (cc CacheClient) GetStruct(ctx context.Context, key string, outputStruct interface{}) error {
	if methodsutil.IsEmpty(key) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := cc.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(serializedValue), &outputStruct); err != nil {
		return err
	}

	return nil
}

func (cc CacheClient) Del(ctx context.Context, keys ...string) error {
	return cc.Redis.Del(ctx, keys...).Err()
}

func (cc CacheClient) DelPattern(ctx context.Context, pattern string) error {
	iter := cc.Redis.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		err := cc.Redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
