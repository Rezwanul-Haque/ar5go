package domain

import "context"

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Get(ctx context.Context, key string) (string, error)
	GetInt(ctx context.Context, key string) (int, error)
	GetStruct(ctx context.Context, key string, outputStruct interface{}) error
	Del(ctx context.Context, keys ...string) error
	DelPattern(ctx context.Context, pattern string) error
}
