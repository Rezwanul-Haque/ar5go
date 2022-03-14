package db

import (
	"ar5go/app/domain"
	"ar5go/infra/logger"
)

var client DatabaseClient

func NewDbClient(lc logger.LogClient) domain.IDb {
	connectMySQL(lc)

	return &DatabaseClient{}
}

func Client() DatabaseClient {
	return client
}
