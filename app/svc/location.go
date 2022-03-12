package svc

import (
	"ar5go/app/domain"
	"ar5go/app/serializers"
	"ar5go/infra/errors"
)

type ILocation interface {
	Create(req serializers.LocationHistoryReq) *errors.RestErr
	Update(req serializers.LocationHistoryReq) (*domain.LocationHistory, *errors.RestErr)
	GetLocationsByUserID(userID uint, filters *serializers.ListFilters) (*serializers.ListFilters, *errors.RestErr)
}
