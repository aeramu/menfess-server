package resolver

import (
	"context"

	auth "github.com/aeramu/menfess-server/internal/auth/service"
	log "github.com/sirupsen/logrus"
)

func (r *resolver) Register(req RegisterReq) (string, error) {
	jwt, err := r.auth.Register(auth.RegisterReq{
		Email:     req.Email,
		Password:  req.Password,
		PushToken: req.PushToken,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"req": req,
		}).Errorln("[Register] Failed register to auth service")
		return "", err
	}
	return jwt, nil
}

func (r *resolver) Login(req LoginReq) (string, error) {
	jwt, err := r.auth.Login(auth.LoginReq{
		Email:     req.Email,
		Password:  req.Password,
		PushToken: req.PushToken,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"req": req,
		}).Errorln("[Login] Failed login on auth service")
		return "", err
	}
	return jwt, nil
}

func (r *resolver) Logout(ctx context.Context) (string, error) {
	jwt := ctx.Value("Authorization").(string)
	if err := r.auth.Logout(auth.LogoutReq{
		Token: jwt,
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorln("[Logout] Failed logout on auth service")
		return "", err
	}
	return "success", nil
}
