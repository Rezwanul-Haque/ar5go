package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidEmail              = NewError("invalid email")
	ErrInvalidPassword           = NewError("invalid password")
	ErrUserRolePermissions       = NewError("failed to fetch role permissions")
	ErrCreateJwt                 = NewError("failed to create JWT token")
	ErrAccessTokenSign           = NewError("failed to sign access_token")
	ErrRefreshTokenSign          = NewError("failed to sign refresh_token")
	ErrStoreTokenUuid            = NewError("failed to store token uuid")
	ErrUpdateLastLogin           = NewError("failed to update last login")
	ErrNoContextUser             = NewError("failed to get user from context")
	ErrInvalidRefreshToken       = NewError("invalid refresh_token")
	ErrInvalidAccessToken        = NewError("invalid access_token")
	ErrInvalidPasswordResetToken = NewError("invalid reset_token")
	ErrInvalidRefreshUuid        = NewError("invalid refresh_uuid")
	ErrInvalidAccessUuid         = NewError("invalid refresh_uuid")
	ErrInvalidJwtSigningMethod   = NewError("invalid signing method while parsing jwt")
	ErrParseJwt                  = NewError("failed to parse JWT token")
	ErrDeleteOldTokenUuid        = NewError("failed to delete old token uuids")
	ErrSendingEmail              = NewError("failed to send email")
	ErrNotAdmin                  = NewError("not admin")
	ErrNotSuperAdmin             = NewError("not super admin")
	ErrEmptyRedisKeyValue        = NewError("empty redis key or value")
	ErrSomethingWentWrong        = "something went wrong"
	ErrRecordNotFound            = "record not found"
	ErrCheckParamBodyHeader      = "check header, params, body"
)

type RestErr struct {
	// example: error message
	Message string `json:"message"`
	// example: 400
	Status int `json:"status"`
	// example: bad_request
	Error string `json:"error"`
}

func NewError(msg string) error {
	return errors.New(msg)
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
