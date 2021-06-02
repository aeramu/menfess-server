package service

import "context"

type FeedbackRepository interface {
	Save(ctx context.Context, feedback Feedback) error
}

type MenfessRequestRepository interface {
	Save(ctx context.Context, request MenfessRequest) error
}
