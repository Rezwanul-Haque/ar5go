package cache

import (
	"ar5go/app/domain"
	"ar5go/infra/logger"
)

var client CacheClient

func NewCacheClient(lc logger.LogClient) domain.ICache {
	connectRedis(lc)

	return &CacheClient{}
}

func Client() CacheClient {
	return client
}
