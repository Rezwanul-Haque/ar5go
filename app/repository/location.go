package repository

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infra/errors"
)

type ILocation interface {
	Save(company *domain.LocationHistory) *errors.RestErr
	Update(*domain.LocationHistory) (*domain.LocationHistory, *errors.RestErr)
	GetLocationsByUserID(userID uint, pagination *serializers.Pagination) ([]*domain.IntermediateLocationHistory, int64, *errors.RestErr)
}
