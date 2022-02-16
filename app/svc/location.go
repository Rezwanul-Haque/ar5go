package svc

import (
	"clean/app/domain"
	"clean/app/serializers"
	"clean/infra/errors"
)

type ILocation interface {
	Create(req serializers.LocationHistoryReq) *errors.RestErr
	Update(req serializers.LocationHistoryReq) (*domain.LocationHistory, *errors.RestErr)
	GetLocationsByUserID(userID uint, pagination *serializers.Pagination) (*serializers.Pagination, *errors.RestErr)
}
