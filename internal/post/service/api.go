package service

type (
	CreateReq struct {
		Body     string
		AuthorID string
		UserID   string
		ParentID string
	}
	GetReq struct {
		ID string
	}
	FeedReq struct {
		First int
		After string
	}
	PostRepliesReq struct {
		PostID string
		First  int
		After  string
	}
	UserPostsReq struct {
		UserID string
		First  int
		After  string
	}
	AuthorPostsReq struct {
		AuthorID string
		First    int
		After    string
	}
	LikeReq struct {
		PostID string
		UserID string
	}
	ReportReq struct {
		PostID string
	}
	DeleteReq struct {
		PostID string
	}
)
