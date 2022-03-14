package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/app/serializers"
	"ar5go/infra/conn/db"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
	"context"
)

type location struct {
	ctx context.Context
	lc  logger.LogClient
	DB  db.DatabaseClient
}

// NewLocationRepository will create an object that represent the Location.Repository implementations
func NewLocationRepository(ctx context.Context, lc logger.LogClient, dbc db.DatabaseClient) repository.ILocation {
	return &location{
		ctx: ctx,
		lc:  lc,
		DB:  dbc,
	}
}

func (r location) SaveLocation(history *domain.LocationHistory) *errors.RestErr {
	return r.DB.SaveLocation(history)
}

func (r location) UpdateLocation(history *domain.LocationHistory) (*domain.LocationHistory, *errors.RestErr) {
	return r.DB.UpdateLocation(history)
}

func (r location) GetLocationsByUserID(userID uint, filters *serializers.ListFilters) ([]*domain.IntermediateLocationHistory, *errors.RestErr) {
	return r.DB.GetLocationsByUserID(userID, filters)
}
