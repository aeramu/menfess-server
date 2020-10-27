package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Service interface{
	Request(req RequestReq) (*Room, error)
	Get(req GetReq) (*Room, error)
	GetList(req GetListReq) (*[]Room, error)
}

func NewService(repo Repository) Service{
	return &service{repo: repo}
}

type service struct{
	repo Repository
}

func (s *service) Request(req RequestReq) (*Room, error){
	room := Room{
		ID:     primitive.NewObjectID().Hex(),
		Name:   req.Name,
		Desc:   req.Desc,
		Avatar: "",
		Status: false,
	}
	if err := s.repo.Save(room); err != nil{
		log.Println("Repository Error:", err)
		return nil, err
	}
	return &room, nil
}

func (s *service) Get(req GetReq) (*Room, error){
	room, err := s.repo.FindByID(req.ID)
	if err != nil{
		log.Println("Repository Error", err)
		return nil, err
	}
	return room, nil
}

func (s *service) GetList(req GetListReq) (*[]Room, error){
	//rooms, err := s.repo.FindByStatus(true)
	rooms, err := s.repo.FindAll()
	if err != nil{
		log.Println("Repository Error", err)
		return nil, err
	}
	return rooms, nil
}