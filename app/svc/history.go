package svc

import (
	"clean/app/serializers"
	"clean/infrastructure/errors"
)

type IHistory interface {
	Create(req serializers.LocationHistoryReq) *errors.RestErr
}
