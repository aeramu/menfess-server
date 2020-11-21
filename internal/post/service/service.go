package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Service interface {
	Create(req CreateReq) (*Post, error)
	Get(req GetReq) (*Post, error)
	Feed(req FeedReq) (*[]Post, error)
	Replies(req RepliesReq) (*[]Post, error)
	Like(req LikeReq) (*Post, error)
	Report(req ReportReq) (*Post, error)
	Delete(req DeleteReq) error
}

func NewService(repo Repository, notif NotificationClient) Service {
	return &service{
		repo: repo,
		notif: notif,
	}
}

type service struct {
	repo  Repository
	notif NotificationClient
}

func (s *service) Create(req CreateReq) (*Post, error) {
	id := primitive.NewObjectID()
	post := Post{
		ID: id.Hex(),
		Timestamp: int(id.Timestamp().Unix()),
		Body: req.Body,
		AuthorID: req.AuthorID,
		UserID: req.UserID,
		ParentID: req.ParentID,
		LikeIDs: map[string]bool{},
		RepliesCount:   0,
		ReportsCount: 0,
	}
	if err := s.repo.Save(post); err != nil{
		log.Println("Repository Error:", err)
		return nil, err
	}

	if post.ParentID == "" {
		return &post, nil
	}

	parent, err := s.repo.FindByID(post.ParentID)
	if err != nil {
		log.Println("Repository Error:", err)
	}
	if parent != nil {
		parent.RepliesCount++
		if err := s.repo.Save(*parent); err != nil{
			log.Println("Repository Error:", err)
		}
		if err := s.notif.Send("reply", parent.UserID, post); err != nil{
			log.Println("Notification Client Error:", err)
		}
	}
	return &post, err
}

func (s *service) Get(req GetReq) (*Post, error) {
	post, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return post, nil
}

func (s *service) Feed(req FeedReq) (*[]Post, error) {
	if req.After == ""{
		req.After = "ffffffffffffffffffffffff"
	}
	postList, err := s.repo.FindByParentID("", req.First, req.After, true)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return postList, nil
}

func (s *service) Replies(req RepliesReq) (*[]Post, error) {
	postList, err := s.repo.FindByParentID(req.PostID, req.First, req.After, false)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return postList, nil
}

func (s *service) Like(req LikeReq) (*Post, error) {
	post, err := s.repo.FindByID(req.PostID)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	if post == nil{
		return nil, nil
	}
	post.Like(req.UserID)
	if err := s.repo.Save(*post); err != nil{
		log.Println("Repository Error:", err)
		return nil, err
	}
	if err := s.notif.Send("like", post.UserID, req); err != nil{
		log.Println("Notification Client Error:", err)
	}
	return post, nil
}

func (s *service) Report(req ReportReq) (*Post, error) {
	panic("implement me")
}

func (s *service) Delete(req DeleteReq) error {
	panic("implement me")
}