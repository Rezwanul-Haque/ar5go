package db

import "ar5go/app/domain"

var client DatabaseClient

func NewDbClient() domain.IDb {
	connectMySQL()

	return &DatabaseClient{}
}

func Client() DatabaseClient {
	return client
}
