package svc

import (
	"ar5go/app/serializers"
)

type IToken interface {
	CreateToken(userID, companyID uint) (*serializers.JwtToken, error)
	StoreTokenUuid(userID, companyID uint, token *serializers.JwtToken) error
	DeleteTokenUuid(uuid ...string) error
}
