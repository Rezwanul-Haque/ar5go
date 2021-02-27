package repository

import (
	"clean/app/domain"
	"clean/infra/errors"
)

type IHistory interface {
	Save(company *domain.LocationHistory) *errors.RestErr
}
