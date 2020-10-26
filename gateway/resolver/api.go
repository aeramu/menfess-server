package resolver

import "github.com/graph-gophers/graphql-go"

type(
	MenfessPostRequest struct{
		ID graphql.ID
	}
	ConnectionRequest struct {
		First *int32
		After *graphql.ID
		Sort  *bool
	}
	MenfessPostRoomsRequest struct {
		IDs   []graphql.ID
		First *int32
		After *graphql.ID
	}
	UpvoteMenfessPostRequest struct {
		PostID graphql.ID
	}
	DownvoteMenfessPostRequest struct {
		PostID graphql.ID
	}
)
