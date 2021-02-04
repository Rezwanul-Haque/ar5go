package serializers

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

type TokenRefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *UserResp `json:"user"`
}

type JwtToken struct {
	UserID        uint   `json:"uid"`
	AccessToken   string `json:"act"`
	RefreshToken  string `json:"rft"`
	AccessUuid    string `json:"aid"`
	RefreshUuid   string `json:"rid"`
	AccessExpiry  int64  `json:"axp"`
	RefreshExpiry int64  `json:"rxp"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (c ChangePasswordReq) Validate() error {
	return v.ValidateStruct(&c,
		v.Field(&c.OldPassword, v.Required),
		v.Field(&c.NewPassword, v.Required, v.Length(8, 0)),
	)
}

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

func (f ForgotPasswordReq) Validate() error {
	return v.ValidateStruct(&f,
		v.Field(&f.Email, v.Required, is.EmailFormat),
	)
}

type VerifyResetPasswordReq struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

func (vr VerifyResetPasswordReq) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.Token, v.Required),
		v.Field(&vr.ID, v.Required),
	)
}

type ResetPasswordReq struct {
	ID       int    `json:"id"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (rp ResetPasswordReq) Validate() error {
	return v.ValidateStruct(&rp,
		v.Field(&rp.Token, v.Required),
		v.Field(&rp.ID, v.Required),
		v.Field(&rp.Password, v.Required, v.Length(8, 0)),
	)
}
