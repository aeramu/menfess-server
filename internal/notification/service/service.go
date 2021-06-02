package service

import (
	"context"

	"github.com/aeramu/menfess-server/internal/notification/constants"

	log "github.com/sirupsen/logrus"
)

//go:generate mockery --all --inpackage --case=underscore
type Service interface {
	AddPushToken(ctx context.Context, req AddPushTokenReq) error
	RemovePushToken(ctx context.Context, req RemovePushTokenReq) error
	SendNotification(ctx context.Context, req SendNotificationReq) error
}

func NewService(repository Repository, pushService PushServiceClient) Service {
	return &service{
		repo: repository,
		push: pushService,
	}
}

type service struct {
	repo Repository
	push PushServiceClient
}

func (s *service) SendNotification(ctx context.Context, req SendNotificationReq) error {
	user, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.UserID,
		}).Errorln("[AddPushToken] Failed get user by id")
		return err
	}
	if user == nil {
		return constants.ErrUserNotFound
	}

	if err := s.push.Send(ctx, user.PushToken, req.Title, req.Body, req.Data); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"pushToken": user.PushToken,
		}).Errorln("[Send Notification] Failed send notification")
		return err
	}

	return nil
}

func (s *service) AddPushToken(ctx context.Context, req AddPushTokenReq) error {
	user, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.ID,
		}).Errorln("[AddPushToken] Failed get user by id")
		return err
	}
	if user == nil {
		return constants.ErrUserNotFound
	}
	user.AddPushToken(req.PushToken)
	if err := s.repo.Save(ctx, *user); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        user.ID,
			"pushToken": user.PushToken,
		}).Errorln("[AddPushToken] Failed save updated user")
		return err
	}
	return nil
}

func (s *service) RemovePushToken(ctx context.Context, req RemovePushTokenReq) error {
	user, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.ID,
		}).Errorln("[RemovePushToken] Failed get user by id")
		return err
	}
	if user == nil {
		return constants.ErrUserNotFound
	}
	user.RemovePushToken(req.PushToken)
	if err := s.repo.Save(ctx, *user); err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"id":        user.ID,
			"pushToken": user.PushToken,
		}).Errorln("[RemovePushToken] Failed save updated user")
		return err
	}
	return err
}
