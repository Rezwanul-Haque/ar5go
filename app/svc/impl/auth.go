package impl

import (
	"boilerplate/app/domain"
	"boilerplate/app/repository"
	"boilerplate/app/serializers"
	"boilerplate/app/svc"
	"boilerplate/app/utils/consts"
	"boilerplate/app/utils/methodutil"
	"boilerplate/app/utils/msgutil"
	"boilerplate/infra/config"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	urepo repository.IUsers
	tSvc  svc.IToken
}

func NewAuthService(urepo repository.IUsers, tokenSvc svc.IToken) svc.IAuth {
	return &auth{
		urepo: urepo,
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

	if token, err = as.tSvc.CreateToken(user.ID); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	if err = as.urepo.SetLastLoginAt(user); err != nil {
		logger.Error("error occur when trying to set last login", err)
		return nil, errors.ErrUpdateLastLogin
	}

	var userResp *serializers.UserWithParamsResp

	if userResp, err = as.getUserInfoWithParam(user.ID); err != nil {
		return nil, err
	}

	res := &serializers.LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         userResp,
	}
	return res, nil
}

func (as *auth) RefreshToken(refreshToken string) (*serializers.LoginResp, error) {
	oldToken, err := as.parseToken(refreshToken, consts.RefreshTokenType)
	if err != nil {
		return nil, errors.ErrInvalidRefreshToken
	}

	var newToken *serializers.JwtToken

	if newToken, err = as.tSvc.CreateToken(oldToken.UserID); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	var userResp *serializers.UserWithParamsResp

	if userResp, err = as.getUserInfoWithParam(newToken.UserID); err != nil {
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

	var resp *serializers.VerifyTokenResp

	if resp, err = as.getTokenResponse(token); err != nil {
		return nil, err
	}

	return resp, nil
}

func (as *auth) getUserInfoWithParam(userID uint) (*serializers.UserWithParamsResp, error) {
	userWithParams := serializers.UserWithParamsResp{}

	var err error

	user, getErr := as.urepo.GetUserWithPermissions(userID, true)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	err = methodutil.StructToStruct(user, &userWithParams)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	return &userWithParams, nil
}

func (as *auth) parseToken(token, tokenType string) (*serializers.JwtToken, error) {
	claims, err := as.parseTokenClaim(token, tokenType)
	if err != nil {
		return nil, err
	}

	tokenDetails := &serializers.JwtToken{}

	if err := methodutil.MapToStruct(claims, &tokenDetails); err != nil {
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

	parsedToken, err := methodutil.ParseJwtToken(token, secret)
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

	logger.Error(err.Error(), err)

	user, getErr := as.urepo.GetTokenUser(token.UserID)
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	err = methodutil.StructToStruct(user, &resp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user to verify token response"), err)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	return resp, err
}
