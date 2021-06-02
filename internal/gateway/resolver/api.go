package resolver

import "github.com/graph-gophers/graphql-go"

type (
	// Auth
	RegisterReq struct {
		Email     string
		Password  string
		PushToken string
	}
	LoginReq struct {
		Email     string
		Password  string
		PushToken string
	}

	// User
	UpdateProfileReq struct {
		Name   *string
		Avatar *string
		Bio    *string
	}

	PostReq struct {
		ID graphql.ID
	}
	ConnectionReq struct {
		First *int32
		After *graphql.ID
	}
	CreatePostReq struct {
		Body     string
		AuthorID graphql.ID
		ParentID *graphql.ID
	}
	DeletePostReq struct {
		ID graphql.ID
	}
	LikePostReq struct {
		ID graphql.ID
	}
)
