package repository

import (
	"clean/app/domain"
	"clean/infrastructure/errors"
)

type IHistory interface {
	Save(company *domain.LocationHistory) *errors.RestErr
}
