package service

import (
	"github.com/aeramu/menfess-server/internal/auth/constants"
)

type (
	RegisterReq struct {
		Email     string
		Password  string
		PushToken string
	}
	LoginReq struct {
		Email     string
		Password  string
		PushToken string
	}
	LogoutReq struct {
		Token string
	}
	AuthReq struct {
		Token string
	}
)

func (req RegisterReq) Validate() error {
	if req.Email == "" {
		return constants.ErrEmailRequied
	}
	// TODO: validate email format
	if req.Email == "" {
		return constants.ErrEmailNotValid
	}
	if req.Password == "" {
		return constants.ErrPasswordRequired
	}
	if req.PushToken == "" {
		return constants.ErrPushTokenRequired
	}
	return nil
}

func (req LoginReq) Validate() error {
	if req.Email == "" {
		return constants.ErrEmailRequied
	}
	// TODO: validate email format
	if req.Email == "" {
		return constants.ErrEmailNotValid
	}
	if req.Password == "" {
		return constants.ErrPasswordRequired
	}
	if req.PushToken == "" {
		return constants.ErrPushTokenRequired
	}
	return nil
}
