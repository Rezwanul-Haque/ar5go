package svc

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/consts"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/config"
	"clean/infra/errors"
	"clean/infra/logger"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	urepo repository.IUsers
	rSvc  svc.ICache
	tSvc  svc.IToken
}

func NewAuthService(urepo repository.IUsers, rSvc svc.ICache, tokenSvc svc.IToken) svc.IAuth {
	return &auth{
		urepo: urepo,
		rSvc:  rSvc,
		tSvc:  tokenSvc,
	}
}

func (as *auth) Login(req *serializers.LoginReq) (*serializers.LoginResp, error) {
	var user *domain.User
	var err error

	if user, err = as.urepo.GetUserByEmail(req.Email); err != nil {
		return nil, errors.ErrInvalidEmail
	}

	if req.Admin && !as.urepo.HasRole(user.ID, consts.RoleIDAdmin) {
		return nil, errors.ErrNotAdmin
	}

	loginPass := []byte(req.Password)
	hashedPass := []byte(*user.Password)

	if err = bcrypt.CompareHashAndPassword(hashedPass, loginPass); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrInvalidPassword
	}

	var token *serializers.JwtToken

	if token, err = as.tSvc.CreateToken(user.ID, user.CompanyID); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	if err = as.tSvc.StoreTokenUuid(user.ID, user.CompanyID, token); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrStoreTokenUuid
	}

	if err = as.urepo.SetLastLoginAt(user); err != nil {
		logger.Error("error occur when trying to set last login", err)
		return nil, errors.ErrUpdateLastLogin
	}

	var userResp *serializers.UserWithParamsResp

	if userResp, err = as.getUserInfoWithParam(user.ID, user.CompanyID, false); err != nil {
		return nil, err
	}

	res := &serializers.LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         userResp,
	}
	return res, nil
}

func (as *auth) Logout(user *serializers.LoggedInUser) error {
	return as.tSvc.DeleteTokenUuid(
		config.Redis().AccessUuidPrefix+user.AccessUuid,
		config.Redis().RefreshUuidPrefix+user.RefreshUuid,
	)
}

func (as *auth) RefreshToken(refreshToken string) (*serializers.LoginResp, error) {
	oldToken, err := as.parseToken(refreshToken, consts.RefreshTokenType)
	if err != nil {
		return nil, errors.ErrInvalidRefreshToken
	}

	if !as.userBelongsToTokenUuid(int(oldToken.UserID), oldToken.RefreshUuid, consts.RefreshTokenType) {
		return nil, errors.ErrInvalidRefreshToken
	}

	var newToken *serializers.JwtToken

	if newToken, err = as.tSvc.CreateToken(oldToken.UserID, oldToken.CompanyID); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	if err = as.tSvc.DeleteTokenUuid(
		config.Redis().AccessUuidPrefix+oldToken.AccessUuid,
		config.Redis().RefreshUuidPrefix+oldToken.RefreshUuid,
	); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrDeleteOldTokenUuid
	}

	if err = as.tSvc.StoreTokenUuid(newToken.UserID, newToken.CompanyID, newToken); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrStoreTokenUuid
	}

	var userResp *serializers.UserWithParamsResp

	if userResp, err = as.getUserInfoWithParam(newToken.UserID, newToken.CompanyID, false); err != nil {
		return nil, err
	}

	res := &serializers.LoginResp{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		User:         userResp,
	}

	return res, nil
}

func (as *auth) VerifyToken(accessToken string) (*serializers.VerifyTokenResp, error) {
	token, err := as.parseToken(accessToken, consts.AccessTokenType)
	if err != nil {
		return nil, errors.ErrInvalidAccessToken
	}

	if !as.userBelongsToTokenUuid(int(token.UserID), token.AccessUuid, consts.AccessTokenType) {
		return nil, errors.ErrInvalidAccessToken
	}

	var resp *serializers.VerifyTokenResp

	if resp, err = as.getTokenResponse(token); err != nil {
		return nil, err
	}

	return resp, nil
}

func (as *auth) getUserInfoWithParam(userID, companyID uint, checkInCache bool) (*serializers.UserWithParamsResp, error) {
	userResp := &serializers.UserResp{}
	userWithParams := serializers.UserWithParamsResp{}

	userCacheKey := config.Redis().UserPrefix + strconv.Itoa(int(userID)) + strconv.Itoa(int(companyID))
	var err error

	if checkInCache {
		if err = as.rSvc.GetStruct(userCacheKey, &userResp); err == nil {
			logger.Info("User served from cache")
			return nil, nil
		}

		logger.Error(err.Error(), err)
	}

	user, getErr := as.urepo.GetUserWithPermissions(userID, true)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	err = methodsutil.StructToStruct(user, &userWithParams)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	if err := as.rSvc.Set(userCacheKey, userWithParams, 0); err != nil {
		logger.Error("setting user data on redis key", err)
	}

	return &userWithParams, nil
}

func (as *auth) parseToken(token, tokenType string) (*serializers.JwtToken, error) {
	claims, err := as.parseTokenClaim(token, tokenType)
	if err != nil {
		return nil, err
	}

	tokenDetails := &serializers.JwtToken{}

	if err := methodsutil.MapToStruct(claims, &tokenDetails); err != nil {
		logger.Error(err.Error(), err)
		return nil, err
	}

	if tokenDetails.UserID == 0 || tokenDetails.AccessUuid == "" || tokenDetails.RefreshUuid == "" {
		logger.Info(fmt.Sprintf("%v", claims))
		return nil, errors.ErrInvalidRefreshToken
	}

	return tokenDetails, nil
}

func (as *auth) parseTokenClaim(token, tokenType string) (jwt.MapClaims, error) {
	secret := config.Jwt().AccessTokenSecret

	if tokenType == consts.RefreshTokenType {
		secret = config.Jwt().RefreshTokenSecret
	}

	parsedToken, err := methodsutil.ParseJwtToken(token, secret)
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok || !parsedToken.Valid {
		return nil, errors.ErrInvalidAccessToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrInvalidAccessToken
	}

	return claims, nil
}

func (as *auth) getTokenResponse(token *serializers.JwtToken) (*serializers.VerifyTokenResp, error) {
	var resp *serializers.VerifyTokenResp
	var err error
	tokenCacheKey := config.Redis().TokenPrefix + strconv.Itoa(int(token.UserID))

	if err = as.rSvc.GetStruct(tokenCacheKey, &resp); err == nil {
		logger.Info("Token user served from cache")
		return resp, nil
	}

	logger.Error(err.Error(), err)

	user, getErr := as.urepo.GetTokenUser(token.UserID)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	err = methodsutil.StructToStruct(user, &resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user to verify token response"), err)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	if err := as.rSvc.Set(tokenCacheKey, resp, 0); err != nil {
		logger.Error("setting user data on redis key", err)
	}

	return resp, err
}

func (as *auth) userBelongsToTokenUuid(userID int, uuid, uuidType string) bool {
	prefix := config.Redis().AccessUuidPrefix

	if uuidType == consts.RefreshTokenType {
		prefix = config.Redis().RefreshUuidPrefix
	}

	redisKey := prefix + uuid

	redisUserId, err := as.rSvc.GetInt(redisKey)
	if err != nil {
		switch err {
		case redis.Nil:
			logger.Error(redisKey, errors.NewError(" not found in redis"))
		default:
			logger.Error(err.Error(), err)
		}
		return false
	}

	if userID != redisUserId {
		return false
	}

	return true
}
