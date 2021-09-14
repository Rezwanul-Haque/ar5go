package svc

import (
	"boilerplate/app/serializers"
)

type IAuth interface {
	Login(req *serializers.LoginReq) (*serializers.LoginResp, error)
	RefreshToken(refreshToken string) (*serializers.LoginResp, error)
	VerifyToken(accessToken string) (*serializers.VerifyTokenResp, error)
}
