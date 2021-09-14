package impl

import (
	"boilerplate/app/repository"
	"boilerplate/app/serializers"
	"boilerplate/app/svc"
	"boilerplate/infra/config"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type token struct {
	urepo repository.IUsers
}

func NewTokenService(urepo repository.IUsers) svc.IToken {
	return &token{
		urepo: urepo,
	}
}

func (t *token) CreateToken(userID uint) (*serializers.JwtToken, error) {
	var err error
	jwtConf := config.Jwt()
	token := &serializers.JwtToken{}

	token.UserID = userID
	token.AccessExpiry = time.Now().Add(time.Minute * jwtConf.AccessTokenExpiry).Unix()
	token.AccessUuid = uuid.New().String()

	token.RefreshExpiry = time.Now().Add(time.Minute * jwtConf.RefreshTokenExpiry).Unix()
	token.RefreshUuid = uuid.New().String()

	user, getErr := t.urepo.GetUser(userID, true)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	atClaims := jwt.MapClaims{}
	atClaims["uid"] = user.ID
	atClaims["aid"] = token.AccessUuid
	atClaims["rid"] = token.RefreshUuid
	atClaims["exp"] = token.AccessExpiry

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token.AccessToken, err = at.SignedString([]byte(jwtConf.AccessTokenSecret))
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrAccessTokenSign
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["uid"] = user.ID
	rtClaims["aid"] = token.AccessUuid
	rtClaims["rid"] = token.RefreshUuid
	rtClaims["exp"] = token.RefreshExpiry

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(jwtConf.RefreshTokenSecret))
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrRefreshTokenSign
	}

	return token, nil
}
