package errors

import (
	"errors"
	"net/http"
)

var (
	ErrDuplicateEntry            uint16 = 1062
	ErrInvalidEmail                     = NewError("invalid email")
	ErrInvalidPhone                     = NewError("invalid phone no")
	ErrInvalidPassword                  = NewError("invalid password")
	ErrUserRolePermissions              = NewError("failed to fetch role permissions")
	ErrCreateJwt                        = NewError("failed to create JWT token")
	ErrAccessTokenSign                  = NewError("failed to sign access_token")
	ErrRefreshTokenSign                 = NewError("failed to sign refresh_token")
	ErrStoreTokenUuid                   = NewError("failed to store token uuid")
	ErrUpdateLastLogin                  = NewError("failed to update last login")
	ErrNoContextUser                    = NewError("failed to get user from context")
	ErrInvalidRefreshToken              = NewError("invalid refresh_token")
	ErrInvalidAccessToken               = NewError("invalid access_token")
	ErrInvalidPasswordResetToken        = NewError("invalid reset_token")
	ErrInvalidConfirmationlToken        = NewError("invalid confirmationl_token")
	ErrInvalidRefreshUuid               = NewError("invalid refresh_uuid")
	ErrInvalidAccessUuid                = NewError("invalid refresh_uuid")
	ErrInvalidJwtSigningMethod          = NewError("invalid signing method while parsing jwt")
	ErrParseJwt                         = NewError("failed to parse JWT token")
	ErrDeleteOldTokenUuid               = NewError("failed to delete old token uuids")
	ErrSendingEmail                     = NewError("failed to send email")
	ErrNotAdmin                         = NewError("not admin")
	ErrNotSuperAdmin                    = NewError("not super admin")
	ErrEmptyRedisKeyValue               = NewError("empty redis key or value")
	ErrPhoneOrEmailExists               = "email or phone number already exists"
	ErrProductNameExists                = "product name already exists"
	ErrUserNameNotUnique                = "username already exists!, try another"
	ErrCompanyNameNotUnique             = "company name already exists!, try another"
	ErrPhoneNoIsUnique                  = "phone no already exists!, try another"
	ErrEmailIsUnique                    = "email already exists!, try another"
	ErrPhoneNameNotUnique               = "phone number already exists!, try another"
	ErrSomethingWentWrong               = "something went wrong"
	ErrRecordNotFound                   = "record not found"
	ErrRecordNotvalid                   = "invalid parameters, check name, email or username "
	ErrCheckParamBodyHeader             = "check header, params, body"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
}

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func NewAlreadyExistError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusConflict,
		Error:   "record_already_exists",
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}
