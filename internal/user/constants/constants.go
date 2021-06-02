package constants

import "github.com/pkg/errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrUserNotFound   = errors.New("user not found")
)
