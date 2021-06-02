package service

type Feedback struct {
	ID       string
	Feedback string
	Marked   bool
	UserID   string
}

type MenfessRequest struct {
	ID     string
	Name   string
	Desc   string
	Marked bool
	UserID string
}
