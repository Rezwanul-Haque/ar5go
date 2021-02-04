package svc

import (
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/infrastructure/config"
	"clean/infrastructure/errors"
	"clean/infrastructure/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type token struct {
	urepo repository.IUsers
	rrepo repository.ICache
}

func NewTokenService(urepo repository.IUsers, rrepo repository.ICache) svc.IToken {
	return &token{
		urepo: urepo,
		rrepo: rrepo,
	}
}

func (t *token) CreateToken(userID uint) (*serializers.JwtToken, error) {
	jwtConf := config.Jwt()
	token := &serializers.JwtToken{}

	token.UserID = userID
	token.AccessExpiry = time.Now().Add(time.Minute * jwtConf.AccessTokenExpiry).Unix()
	token.AccessUuid = uuid.New().String()

	token.RefreshExpiry = time.Now().Add(time.Minute * jwtConf.RefreshTokenExpiry).Unix()
	token.RefreshUuid = uuid.New().String()

	user, err := t.urepo.GetUser(userID)
	if err != nil {
		return nil, err
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

func (t *token) StoreTokenUuid(userID uint, token *serializers.JwtToken) error {
	now := time.Now().Unix()

	err := t.rrepo.Set(
		config.Redis().AccessUuidPrefix + token.AccessUuid,
		userID, int(token.AccessExpiry - now),
	)
	if err != nil {
		return err
	}

	err = t.rrepo.Set(
		config.Redis().RefreshUuidPrefix + token.RefreshUuid,
		userID, int(token.RefreshExpiry - now),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *token) DeleteTokenUuid(uuid ...string) error {
	return t.rrepo.Del(uuid...)
}
