package service

type CreateFeedbackReq struct {
	UserID   string
	Feedback string
}

type CreateMenfessRequestReq struct {
	UserID string
	Name   string
	Desc   string
}
