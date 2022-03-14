package impl

import (
	"ar5go/app/repository"
	"ar5go/app/serializers"
	"ar5go/app/svc"
	"ar5go/infra/config"
	"ar5go/infra/conn/cache"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
	"context"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type token struct {
	ctx   context.Context
	lc    logger.LogClient
	urepo repository.IUsers
}

func NewTokenService(ctx context.Context, lc logger.LogClient, urepo repository.IUsers) svc.IToken {
	return &token{
		ctx:   ctx,
		lc:    lc,
		urepo: urepo,
	}
}

func (t *token) CreateToken(userID, companyID uint) (*serializers.JwtToken, error) {
	var err error
	jwtConf := config.Jwt()
	token := &serializers.JwtToken{}

	token.UserID = userID
	token.CompanyID = companyID
	token.AccessExpiry = time.Now().Add(time.Minute * jwtConf.AccessTokenExpiry).Unix()
	token.AccessUuid = uuid.New().String()

	token.RefreshExpiry = time.Now().Add(time.Minute * jwtConf.RefreshTokenExpiry).Unix()
	token.RefreshUuid = uuid.New().String()

	user, getErr := t.urepo.GetUserWithPermissions(userID, true)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	atClaims := jwt.MapClaims{}
	atClaims["uid"] = user.ID
	atClaims["cid"] = user.CompanyID
	atClaims["aid"] = token.AccessUuid
	atClaims["rid"] = token.RefreshUuid
	atClaims["exp"] = token.AccessExpiry

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token.AccessToken, err = at.SignedString([]byte(jwtConf.AccessTokenSecret))
	if err != nil {
		t.lc.Error(err.Error(), err)
		return nil, errors.ErrAccessTokenSign
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["uid"] = user.ID
	rtClaims["cid"] = user.CompanyID
	rtClaims["aid"] = token.AccessUuid
	rtClaims["rid"] = token.RefreshUuid
	rtClaims["exp"] = token.RefreshExpiry

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(jwtConf.RefreshTokenSecret))
	if err != nil {
		t.lc.Error(err.Error(), err)
		return nil, errors.ErrRefreshTokenSign
	}

	return token, nil
}

func (t *token) StoreTokenUuid(userID, companyID uint, token *serializers.JwtToken) error {
	now := time.Now().Unix()
	key, _ := strconv.Atoi(strconv.Itoa(int(userID)) + strconv.Itoa(int(companyID)))

	err := cache.Client().Set(
		t.ctx,
		config.Cache().Redis.AccessUuidPrefix+token.AccessUuid,
		key, int(token.AccessExpiry-now),
	)
	if err != nil {
		return err
	}

	err = cache.Client().Set(
		t.ctx,
		config.Cache().Redis.RefreshUuidPrefix+token.RefreshUuid,
		key, int(token.RefreshExpiry-now),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *token) DeleteTokenUuid(uuid ...string) error {
	return cache.Client().Del(t.ctx, uuid...)
}
