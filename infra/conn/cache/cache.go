package cache

import (
	"ar5go/app/utils/methodsutil"
	"ar5go/infra/errors"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

const KeyPrefix = "ar5go:"

func (rc CacheClient) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	if methodsutil.IsEmpty(key) || methodsutil.IsEmpty(value) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rc.Redis.Set(ctx, KeyPrefix+key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
}

func (rc CacheClient) Get(ctx context.Context, key string) (string, error) {
	if methodsutil.IsEmpty(key) {
		return "", errors.ErrEmptyRedisKeyValue
	}

	return rc.Redis.Get(ctx, key).Result()
}

func (rc CacheClient) GetInt(ctx context.Context, key string) (int, error) {
	if methodsutil.IsEmpty(key) {
		return 0, errors.ErrEmptyRedisKeyValue
	}

	str, err := rc.Redis.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func (rc CacheClient) GetStruct(ctx context.Context, key string, outputStruct interface{}) error {
	if methodsutil.IsEmpty(key) {
		return errors.ErrEmptyRedisKeyValue
	}

	serializedValue, err := rc.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(serializedValue), &outputStruct); err != nil {
		return err
	}

	return nil
}

func (rc CacheClient) Del(ctx context.Context, keys ...string) error {
	return rc.Redis.Del(ctx, keys...).Err()
}

func (rc CacheClient) DelPattern(ctx context.Context, pattern string) error {
	iter := rc.Redis.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		err := rc.Redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
