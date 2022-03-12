package svc

import "ar5go/app/serializers"

type ISystem interface {
	GetHealth() (*serializers.HealthResp, error)
}
