package svc

import (
	"clean/app/serializers"
	"clean/infra/errors"
)

type IHistory interface {
	Create(req serializers.LocationHistoryReq) *errors.RestErr
}
