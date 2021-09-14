package impl

import (
	"boilerplate/app/domain"
	"boilerplate/app/repository"
	"boilerplate/app/serializers"
	"boilerplate/app/svc"
	"boilerplate/app/utils/methodutil"
	"boilerplate/app/utils/msgutil"
	"boilerplate/infra/config"
	"boilerplate/infra/errors"
	"boilerplate/infra/logger"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type users struct {
	urepo repository.IUsers
	msvc  svc.IMails
}

func NewUsersService(urepo repository.IUsers, msvc svc.IMails) svc.IUsers {
	return &users{
		urepo: urepo,
		msvc:  msvc,
	}
}

func (u *users) UserNameIsUnique(req string) error {
	err := u.urepo.UserNameIsUnique(req)
	if err != nil {
		return err
	}
	return nil
}

func (u *users) EmailIsUnique(req *serializers.EmailIsUnique) error {
	err := u.urepo.EmailIsUnique(req.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *users) CreateUser(user domain.User) (*domain.User, *errors.RestErr) {

	hass_pass, _ := GenerateUserCreateRequestPassword()
	user.Password = hass_pass
	resp, saveErr := u.urepo.Save(&user)
	if saveErr != nil {
		return nil, saveErr
	}

	return resp, nil
}

func (u *users) GetUserById(userId uint) (*domain.User, *errors.RestErr) {
	resp, getErr := u.urepo.GetUserByID(userId)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}

func (u *users) GetUserByEmail(userName string) (*domain.User, error) {
	resp, getErr := u.urepo.GetUserByEmail(userName)
	if getErr != nil {
		return nil, getErr
	}
	return resp, nil
}

func (u *users) UpdateUser(userID uint, req serializers.UserReq) *errors.RestErr {
	var user domain.User

	err := methodutil.StructToStruct(req, &user)
	if err != nil {
		logger.ErrorAsJson(msgutil.EntityStructToStructFailedMsg("update user"), err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	user.ID = userID

	if updateErr := u.urepo.Update(&user); updateErr != nil {
		return updateErr
	}
	return nil
}

func (u *users) ChangePassword(id int, data *serializers.ChangePasswordReq) error {
	if payloadErr := data.Validate(); payloadErr != nil {
		return payloadErr
	}
	user, getErr := u.urepo.GetUserByID(uint(id))
	if getErr != nil {
		return errors.NewError(getErr.Message)
	}

	currentPass := []byte(*user.Password)
	if err := bcrypt.CompareHashAndPassword(currentPass, []byte(data.OldPassword)); err != nil {
		logger.ErrorAsJson(msgutil.EntityGenericFailedMsg("comparing hash and old password"), err)
		return errors.ErrInvalidPassword
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 8)

	updates := map[string]interface{}{
		"password":    hashedPass,
		"first_login": false,
	}

	upErr := u.urepo.UpdatePassword(user.ID, updates)
	if upErr != nil {
		return errors.NewError(upErr.Message)
	}

	return nil
}

func (u *users) VerifyResetPassword(req *serializers.VerifyResetPasswordReq) error {
	user, getErr := u.urepo.GetUserByID(uint(req.ID))
	if getErr != nil {
		return errors.NewError(getErr.Message)
	}

	secret := passwordResetSecret(user)

	parsedToken, err := methodutil.ParseJwtToken(req.Token, secret)
	if err != nil {
		logger.ErrorAsJson("error occur when parse jwt token with secret", err)
		return errors.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok && !parsedToken.Valid {
		return errors.ErrInvalidPasswordResetToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.ErrInvalidPasswordResetToken
	}

	parsedEmail := claims["email"].(string)
	if user.Email != parsedEmail {
		return errors.ErrInvalidPasswordResetToken
	}

	return nil
}

func (u *users) ResetPassword(req *serializers.ResetPasswordReq) error {
	if payloadErr := req.Validate(); payloadErr != nil {
		return payloadErr
	}
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 8)

	if err := u.urepo.ResetPassword(req.ID, hashedPass); err != nil {
		return err
	}

	return nil
}

func (u *users) ForgotPassword(email string) error {
	user, err := u.urepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	secret := passwordResetSecret(user)

	payload := jwt.MapClaims{}
	payload["email"] = user.Email

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		logger.ErrorAsJson("error occur when getting complete signed token", err)
		return err
	}

	fpassReq := &serializers.ForgetPasswordMailReq{
		To:     user.Email,
		UserID: user.ID,
		Token:  signedToken,
	}

	logger.InfoAsJson("email payload", fpassReq)

	if config.App().SendEmail {
		if err := u.msvc.SendForgotPasswordEmail(*fpassReq); err != nil {
			return errors.ErrSendingEmail
		}
	}

	return nil
}

func passwordResetSecret(user *domain.User) string {
	return *user.Password + strconv.Itoa(int(user.CreatedAt.Unix()))
}

func GenerateUserCreateRequestPassword() (*string, string) {
	mockPass := config.App().MockPassword

	var password *string
	password = &mockPass
	if !config.App().MockPasswordEnabled {
		*password = methodutil.GenerateRandomStringOfLength(8)
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(*password), 8)
	hasspassword := string(hashedPass)
	return &hasspassword, *password
}
