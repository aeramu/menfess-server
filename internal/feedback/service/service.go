package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CreateFeedback(ctx context.Context, req CreateFeedbackReq) error
	CreateMenfessRequest(ctx context.Context, req CreateMenfessRequestReq) error
}

func NewService(feedbackRepository FeedbackRepository, menfessRequestRepository MenfessRequestRepository) Service {
	return &service{
		feedback: feedbackRepository,
		menfess:  menfessRequestRepository,
	}
}

type service struct {
	feedback FeedbackRepository
	menfess  MenfessRequestRepository
}

func (s *service) CreateFeedback(ctx context.Context, req CreateFeedbackReq) error {
	feedback := Feedback{
		ID:       primitive.NewObjectID().Hex(),
		Feedback: req.Feedback,
		Marked:   false,
		UserID:   req.UserID,
	}
	if err := s.feedback.Save(ctx, feedback); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"feedback": req.Feedback,
			"userID":   req.UserID,
		}).Errorln("[CreateFeedback] Failed save feedback to repo")
		return err
	}
	return nil
}

func (s *service) CreateMenfessRequest(ctx context.Context, req CreateMenfessRequestReq) error {
	menfess := MenfessRequest{
		ID:     primitive.NewObjectID().Hex(),
		Name:   req.Name,
		Desc:   req.Desc,
		Marked: false,
		UserID: req.UserID,
	}
	if err := s.menfess.Save(ctx, menfess); err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"name":   req.Name,
			"desc":   req.Desc,
			"userID": req.UserID,
		}).Errorln("[CreateFeedback] Failed save feedback to repo")
		return err
	}
	return nil
}
