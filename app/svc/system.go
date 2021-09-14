package svc

import (
	"boilerplate/app/serializers"
)

type ISystem interface {
	GetHealth() (*serializers.HealthResp, error)
}
