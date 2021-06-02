package resolver

import (
	"context"

	auth "github.com/aeramu/menfess-server/internal/auth/service"
	user "github.com/aeramu/menfess-server/internal/user/service"
	log "github.com/sirupsen/logrus"
)

func (r *resolver) Me(ctx context.Context) (*User, error) {
	jwt := ctx.Value("Authorization").(string)
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil {
		log.WithError(err).Errorln("[Me] Failed auth on auth service")
		return nil, err
	}
	u, err := r.user.Get(user.GetReq{ID: payload.ID})
	if err != nil {
		log.Println("User Service Error:", err)
		return nil
	}
	if u == nil {
		return nil
	}
	return &User{*u, r}
}

func (r *resolver) UpdateProfile(ctx context.Context, req UpdateProfileReq) (*User, error) {
	jwt := ctx.Value("Authorization").(string)
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"req": req,
		}).Errorln("[UpdateProfile] Failed auth on auth service")
		return nil, err
	}
	u, err := r.user.UpdateProfile(user.UpdateProfileReq{
		ID:     payload.ID,
		Name:   req.Name,
		Avatar: req.Avatar,
		Bio:    req.Bio,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"req": req,
		}).Errorln("[UpdateProfile] Failed update profile on user service")
		return nil, err
	}

	return &User{User: *u, root: r}, nil
}
