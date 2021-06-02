package constants

import "github.com/pkg/errors"

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrEmailRequied       = errors.New("email is required")
	ErrEmailNotValid      = errors.New("email is not valid")
	ErrEmailNotRegistered = errors.New("email doesn't registered")
	ErrEmailRegistered    = errors.New("email already registered")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordWrong      = errors.New("wrong password")
	ErrPushTokenRequired  = errors.New("push token is required")
	ErrInternalServer     = errors.New("internal server error")
)
