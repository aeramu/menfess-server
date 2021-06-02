package constants

import "github.com/pkg/errors"

var (
	ErrUserNotFound = errors.New("user doesn't exist")
)
