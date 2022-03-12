package cache

import "ar5go/app/domain"

var client CacheClient

func NewCacheClient() domain.ICache {
	connectRedis()

	return &CacheClient{}
}

func Client() CacheClient {
	return client
}
