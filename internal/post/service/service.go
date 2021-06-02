package service

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(req CreateReq) (*Post, error)
	Get(req GetReq) (*Post, error)
	Feed(req FeedReq) ([]Post, error)
	PostReplies(req PostRepliesReq) ([]Post, error)
	UserPosts(req UserPostsReq) ([]Post, error)
	AuthorPosts(req AuthorPostsReq) ([]Post, error)
	Like(req LikeReq) (*Post, error)
	Report(req ReportReq) (*Post, error)
	Delete(req DeleteReq) error
}

func NewService(repo Repository, notif NotificationClient) Service {
	return &service{
		repo:  repo,
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
		ID:           id.Hex(),
		Timestamp:    int(id.Timestamp().Unix()),
		Body:         req.Body,
		AuthorID:     req.AuthorID,
		UserID:       req.UserID,
		ParentID:     req.ParentID,
		LikeIDs:      map[string]bool{},
		RepliesCount: 0,
		ReportsCount: 0,
	}
	if err := s.repo.Save(post); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"id":       id,
			"body":     req.Body,
			"authorID": req.AuthorID,
			"userID":   req.UserID,
			"parentID": req.ParentID,
		}).Errorln("[Create] Failed save post")
		return nil, err
	}

	if post.ParentID == "" {
		return &post, nil
	}

	parent, err := s.repo.FindByID(post.ParentID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  post.ParentID,
		}).Warningln("[Create] Failed get parent post")
		return &post, nil
	}
	if parent == nil {
		log.Warningln("[Create] Parent post not found")
		return &post, nil
	}

	parent.RepliesCount++
	if err := s.repo.Save(*parent); err != nil {
		log.WithFields(log.Fields{
			"err":          err,
			"id":           id,
			"parentID":     req.ParentID,
			"repliesCount": parent.RepliesCount,
		}).Warningln("[Create] Failed save parent post")
	}

	if err := s.notif.Send("reply", parent.UserID, post); err != nil {
		log.WithFields(log.Fields{
			"err":          err,
			"event":        "reply",
			"parentUserID": parent.UserID,
			"post":         post,
		}).Warningln("[Create] Failed send notification")
	}

	return &post, nil
}

func (s *service) Get(req GetReq) (*Post, error) {
	post, err := s.repo.FindByID(req.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.ID,
		}).Errorln("[Get] Failed get post")
		return nil, err
	}
	return post, nil
}

func (s *service) Feed(req FeedReq) ([]Post, error) {
	if req.After == "" {
		req.After = "ffffffffffffffffffffffff"
	}
	postList, err := s.repo.FindByParentID("", req.First, req.After, true)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"first": req.First,
			"after": req.After,
		}).Errorln("[Feed] Failed get posts by parent id")
		return nil, err
	}
	return postList, nil
}

func (s *service) PostReplies(req PostRepliesReq) ([]Post, error) {
	postList, err := s.repo.FindByParentID(req.PostID, req.First, req.After, false)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"parentID": req.PostID,
			"first":    req.First,
			"after":    req.After,
		}).Errorln("[PostReplies] Failed get posts by parent id")
		return nil, err
	}
	return postList, nil
}

func (s *service) UserPosts(req UserPostsReq) ([]Post, error) {
	postList, err := s.repo.FindByUserID(req.UserID, req.First, req.After, true)
	if err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"userID": req.UserID,
			"first":  req.First,
			"after":  req.After,
		}).Errorln("[UserPosts] Failed get posts by parent id")
		return nil, err
	}
	return postList, nil
}

func (s *service) AuthorPosts(req AuthorPostsReq) ([]Post, error) {
	postList, err := s.repo.FindByAuthorID(req.AuthorID, req.First, req.After, true)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"authorID": req.AuthorID,
			"first":    req.First,
			"after":    req.After,
		}).Errorln("[AuthorPosts] Failed get posts by parent id")
		return nil, err
	}
	return postList, nil
}

func (s *service) Like(req LikeReq) (*Post, error) {
	post, err := s.repo.FindByID(req.PostID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.PostID,
		}).Errorln("[Like] Failed get post")
		return nil, err
	}
	if post == nil {
		log.Warningln("[Like] Post not found")
		return nil, nil
	}

	post.Like(req.UserID)
	if err := s.repo.Save(*post); err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"id":     post.ID,
			"userID": req.UserID,
			"like":   post.LikeIDs,
		}).Errorln("[Like] Failed save liked post")
		return nil, err
	}

	if err := s.notif.Send("like", post.UserID, req); err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"event":  "like",
			"userID": post.UserID,
			"req":    req,
		}).Warningln("[Like] Failed send notification")
	}

	return post, nil
}

func (s *service) Report(req ReportReq) (*Post, error) {
	post, err := s.repo.FindByID(req.PostID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  req.PostID,
		}).Errorln("[Report] Failed get post")
		return nil, err
	}
	if post == nil {
		log.Warningln("[Report] Post not found")
		return nil, nil
	}

	post.ReportsCount++
	if err := s.repo.Save(*post); err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"id":  post.ID,
		}).Errorln("[Report] Failed save liked post")
		return nil, err
	}

	return post, nil
}

func (s *service) Delete(req DeleteReq) error {
	panic("implement me")
}
