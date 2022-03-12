package impl

import (
	"ar5go/app/repository"
	"ar5go/infra/conn/cache"
	"ar5go/infra/conn/db"
	"ar5go/infra/logger"
	"context"
	"fmt"
)

type system struct {
	ctx   context.Context
	DB    db.DatabaseClient
	Cache cache.CacheClient
}

// NewSystemRepository will create an object that represent the System.Repository implementations
func NewSystemRepository(ctx context.Context, dbc db.DatabaseClient, c cache.CacheClient) repository.ISystem {
	return &system{
		ctx:   ctx,
		DB:    dbc,
		Cache: c,
	}
}

func (r *system) DBCheck() (bool, error) {
	dB, _ := r.DB.DB.DB()
	if err := dB.Ping(); err != nil {
		return false, err
	}

	return true, nil
}

func (r *system) CacheCheck() bool {
	pong, err := r.Cache.Redis.Ping(r.ctx).Result()
	if err != nil {
		return false
	}

	logger.Info(fmt.Sprintf("%v from cache", pong))

	return true
}
