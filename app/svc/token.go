package svc

import (
	"clean/app/serializers"
)

type IToken interface {
	CreateToken(userID uint) (*serializers.JwtToken, error)
}
