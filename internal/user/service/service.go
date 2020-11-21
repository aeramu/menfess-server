package service

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(req CreateReq) (*User, error)
	Get(req GetReq) (*User, error)
	GetByEmail(req GetByEmailReq) (*User, error)
	GetMenfess(req GetMenfessReq) (*[]User, error)
	UpdateProfile(req UpdateProfileReq) (*User, error)
	AddPushToken(req PushTokenReq) error
	RemovePushToken(req PushTokenReq) error
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
		ID:        primitive.NewObjectID().Hex(),
		Email:     req.Email,
		Password:  req.Password,
		Name:      "",
		Avatar:    "",
		Bio:       "",
		PushToken: map[string]bool{req.PushToken: true},
	}
	if err := s.repo.Save(*user); err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return user, nil
}

func (s *service) Get(req GetReq) (*User, error) {
	user, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return user, nil
}

func (s *service) GetByEmail(req GetByEmailReq) (*User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return user, nil
}

func (s *service) GetMenfess(req GetMenfessReq) (*[]User, error) {
	users, err := s.repo.FindByType("menfess")
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return users, nil
}

func (s *service) UpdateProfile(req UpdateProfileReq) (*User, error) {
	user, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	user.Name = req.Name
	user.Avatar = req.Avatar
	user.Bio = req.Bio
	if err := s.repo.Save(*user); err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return user, nil
}

func (s *service) AddPushToken(req PushTokenReq) error {
	user, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.Println("Repository Error:", err)
		return err
	}
	if user == nil {
		return errors.New("user doesn't exist")
	}
	user.PushToken[req.PushToken] = true
	if err := s.repo.Save(*user); err != nil {
		log.Println("Repository Error:", err)
		return err
	}
	return nil
}

func (s *service) RemovePushToken(req PushTokenReq) error {
	user, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.Println("Repository Error:", err)
		return err
	}
	if user == nil {
		return errors.New("user doesn't exist")
	}
	user.PushToken[req.PushToken] = false
	if err := s.repo.Save(*user); err != nil {
		log.Println("Repository Error:", err)
		return err
	}
	return err
}
