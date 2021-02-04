package svc

import (
	"clean/app/serializers"
)

type IAuth interface {
	Login(req *serializers.LoginReq) (*serializers.LoginResp, error)
	Logout(user *serializers.LoggedInUser) error
	RefreshToken(refreshToken string) (*serializers.LoginResp, error)
	VerifyToken(accessToken string) (*serializers.VerifyTokenResp, error)
}
