package resolver

import (
	"context"
	"log"

	auth "github.com/aeramu/menfess-server/internal/auth/service"
	post "github.com/aeramu/menfess-server/internal/post/service"
	user "github.com/aeramu/menfess-server/internal/user/service"
	"github.com/graph-gophers/graphql-go"
)

type Post struct {
	post.Post
	root *resolver
}

func (r Post) ID() graphql.ID {
	return graphql.ID(r.Post.ID)
}
func (r Post) Timestamp() int32 {
	return int32(r.Post.Timestamp)
}
func (r Post) Body() string {
	return r.Post.Body
}
func (r Post) Author() *User {
	u, err := r.root.user.Get(user.GetReq{ID: r.Post.AuthorID})
	if err != nil {
		log.Println("User Service Error:", err)
		return nil
	}
	return &User{*u, r.root}
}
func (r Post) LikesCount() int32 {
	return int32(r.Post.LikesCount())
}
func (r Post) Liked(ctx context.Context) bool {
	jwt := ctx.Value("Authorization").(string)
	payload, err := r.root.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil {
		log.Println("Auth Service Error:", err)
		return false
	}
	return r.IsLiked(payload.ID)
}
func (r Post) RepliesCount() int32 {
	return int32(r.Post.RepliesCount)
}
func (r Post) Replies(req ConnectionReq) PostConnection {
	first := 20
	if req.First != nil {
		first = int(*req.First)
	}
	after := ""
	if req.After != nil {
		after = string(*req.After)
	}
	postList, err := r.root.post.PostReplies(post.PostRepliesReq{
		PostID: r.Post.ID,
		First:  first,
		After:  after,
	})
	if err != nil {
		log.Println("Post Service Error:", err)
		return PostConnection{[]Post{}, r.root}
	}
	var posts []Post
	for _, elem := range postList {
		posts = append(posts, Post{elem, r.root})
	}
	return PostConnection{posts, r.root}
}

type PostConnection struct {
	posts []Post
	root  *resolver
}

func (r PostConnection) Edges() []Post {
	return r.posts
}
func (r PostConnection) PageInfo() PageInfo {
	var nodeList PageInfo
	for _, node := range r.posts {
		nodeList = append(nodeList, node)
	}
	return nodeList
}

type Node interface{ ID() graphql.ID }
type PageInfo []Node

func (r PageInfo) StartCursor() *graphql.ID {
	if len(r) == 0 {
		return nil
	}
	startCursor := r[0].ID()
	return &startCursor
}
func (r PageInfo) EndCursor() *graphql.ID {
	if len(r) == 0 {
		return nil
	}
	endCursor := r[len(r)-1].ID()
	return &endCursor
}
