package svc

import "clean/app/serializers"

type ISystem interface {
	GetHealth() (*serializers.HealthResp, error)
}
