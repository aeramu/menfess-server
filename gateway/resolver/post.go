package resolver

import (
	post "github.com/aeramu/menfess-server/post/service"
	room "github.com/aeramu/menfess-server/room/service"
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
func (r Post) Name() string{
	return r.Post.Name
}
func (r Post) Body() string{
	return r.Post.Body
}
func (r Post) Avatar() string{
	return r.Post.Avatar
}
func (r Post) Room() string {
	if r.Post.RoomID == "" {
		return "General"
	}
	//return r.Post.RoomID.Name
	room, err := r.root.room.Get(room.GetReq{ID: r.Post.RoomID})
	if err != nil{
		return ""
	}
	return room.Name
}
func (r Post) ReplyCount() int32 {
	return int32(r.Post.ReplyCount)
}
func (r Post) UpvoteCount() int32 {
	return int32(r.Post.UpvoteCount())
}
func (r Post) DownvoteCount() int32 {
	return int32(r.Post.DownvoteCount())
}
func (r Post) Upvoted() bool {
	accountID := r.root.Context.Value("request").(map[string]string)["id"]
	return r.IsUpvoted(accountID)
}
func (r Post) Downvoted() bool {
	accountID := r.root.Context.Value("request").(map[string]string)["id"]
	return r.IsDownvoted(accountID)
}
func (r Post) Parent() *Post {
	if r.Post.ParentID == "" {
		return nil
	}
	p, err := r.root.post.Get(r.ParentID)
	if err != nil{
		return nil
	}
	return &Post{*p, r.root}
}
func (r Post) Repost() *Post {
	if r.Post.RepostID == ""{
		return nil
	}
	p, err := r.root.post.Get(r.RepostID)
	if err != nil{
		return nil
	}
	return &Post{*p, r.root}
}
func (r Post) Child(args ConnectionRequest) PostConnection {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := ""
	if args.After != nil{
		after = string(*args.After)
	}
	postList, err := r.root.post.Child(r.Post.ID, first, after)
	if err != nil {
		return PostConnection{}
	}
	var posts []Post
	for _, elem := range *postList {
		posts = append(posts, Post{elem, r.root})
	}
	return PostConnection{posts, r.root}
}

type PostConnection struct{
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

type Node interface { ID() graphql.ID }
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

