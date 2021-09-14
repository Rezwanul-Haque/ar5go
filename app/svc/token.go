package svc

import (
	"boilerplate/app/serializers"
)

type IToken interface {
	CreateToken(userID uint) (*serializers.JwtToken, error)
}
