package service

import (
	"github.com/aeramu/menfess-server/internal/user/constants"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(req CreateReq) (*User, error)
	Get(req GetReq) (*User, error)
	GetMenfess(req GetMenfessReq) ([]User, error)
	UpdateProfile(req UpdateProfileReq) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(req CreateReq) (*User, error) {
	user := &User{
		ID:     primitive.NewObjectID().Hex(),
		Name:   "",
		Avatar: "",
		Bio:    "",
		Type:   req.Type,
	}
	if err := s.repo.Save(*user); err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  user.ID,
		}).Errorln("[Create] Failed save created user")
		return nil, constants.ErrInternalServer
	}
	return user, nil
}

func (s *service) Get(req GetReq) (*User, error) {
	user, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.ID,
		}).Errorln("[Get] Failed get user by id")
		return nil, constants.ErrInternalServer
	}
	return user, nil
}

func (s *service) GetMenfess(req GetMenfessReq) ([]User, error) {
	users, err := s.repo.FindByType("menfess")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorln("[GetMenfess] Failed get user by type menfess")
		return nil, err
	}
	return users, nil
}

func (s *service) UpdateProfile(req UpdateProfileReq) (*User, error) {
	user, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.ID,
		}).Errorln("[UpdateProfile] Failed get user by id")
		return nil, constants.ErrInternalServer
	}
	if user == nil {
		return nil, constants.ErrUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if err := s.repo.Save(*user); err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"id":     user.ID,
			"name":   user.Name,
			"avatar": user.Avatar,
			"bio":    user.Bio,
		}).Errorln("[Create] Failed save updated user")
		return nil, constants.ErrInternalServer
	}
	return user, nil
}
