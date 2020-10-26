package resolver

import (
	"github.com/aeramu/menfess-server/post/service"
	"github.com/graph-gophers/graphql-go"
)

//Post grahql
type Post interface {
	ID() graphql.ID
	Timestamp() int32
	Name() string
	Avatar() string
	Body() string
	Room() string
	ReplyCount() int32
	UpvoteCount() int32
	DownvoteCount() int32
	Upvoted() bool
	Downvoted() bool
	Parent() Post
	Repost() Post
	Child(args struct {
		First  *int32
		After  *graphql.ID
		Before *graphql.ID
		Sort   *int32
	}) PostConnection
}

// post graphql
type post struct {
	service.Post
	pr *resolver
}

// ID graphql
func (r post) ID() graphql.ID {
	return graphql.ID(r.Post.ID())
}

// Timestamp graphql
func (r post) Timestamp() int32 {
	return int32(r.Post.Timestamp())
}

//Room graphql
func (r post) Room() string {
	if r.Post.Room() == nil {
		return "General"
	}
	return r.Post.Room().Name()
}

// ReplyCount graphql
func (r *post) ReplyCount() int32 {
	return int32(r.Post.ReplyCount())
}

// UpvoteCount graphql
func (r *post) UpvoteCount() int32 {
	return int32(r.Post.UpvoteCount())
}

// DownvoteCount graphql
func (r *post) DownvoteCount() int32 {
	return int32(r.Post.DownvoteCount())
}

//Upvoted bool
func (r *post) Upvoted() bool {
	accountID := r.pr.Context.Value("request").(map[string]string)["id"]
	return r.IsUpvoted(accountID)
}

//Downvoted bool
func (r *post) Downvoted() bool {
	accountID := r.pr.Context.Value("request").(map[string]string)["id"]
	return r.IsDownvoted(accountID)
}

// Parent graphql
func (r *post) Parent() Post {
	return nil
}

//Repost graphql
func (r *post) Repost() Post {
	if r.Post.Repost() == nil {
		return nil
	}
	return &post{r.Post.Repost(), r.pr}
}

// Child graphql
func (r *post) Child(args struct {
	First  *int32
	After  *graphql.ID
	Before *graphql.ID
	Sort   *int32
}) PostConnection {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := "000000000000000000000000"
	postList := r.pr.Interactor.PostChild(r.Post.ID(), first, after)
	return &postConnection{postList, r.pr}
}
