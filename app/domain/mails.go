package domain

type Mail struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Password string `json:"password"`
}

type CompanyCreatedMailReq struct {
	To        string `json:"to" valid:"required"`
	Subject   string `json:"subject"`
	CompanyID uint   `json:"user_id" valid:"required"`
	Password  string `json:"password" valid:"required"`
	Token     string `json:"token"`
}

type UserCreateMail struct {
	To       string `json:"to"`
	UserID   uint   `json:"user_id"`
	Password string `json:"password" valid:"required"`
	Token    string `json:"token"`
}

type ForgetPasswordMail struct {
	To     string `json:"to"`
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

type GenericMail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
