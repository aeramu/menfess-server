package resolver

import (
	"github.com/aeramu/menfess-server/entity"
	"github.com/graph-gophers/graphql-go"
)

// postResolver graphql
type postResolver struct {
	post entity.Post
	pr   *resolver
}

// ID graphql
func (r *postResolver) ID() graphql.ID {
	return graphql.ID(r.post.ID())
}

// Timestamp graphql
func (r *postResolver) Timestamp() int32 {
	return int32(r.post.Timestamp())
}

// Name graphql
func (r *postResolver) Name() string {
	return r.post.Name()
}

// Avatar graphql
func (r *postResolver) Avatar() string {
	return r.post.Avatar()
}

// Body graphql
func (r *postResolver) Body() string {
	return r.post.Body()
}

//Room graphql
func (r *postResolver) Room() string {
	if r.post.Room() == nil {
		return "General"
	}
	return r.post.Room().Name()
}

// ReplyCount graphql
func (r *postResolver) ReplyCount() int32 {
	return int32(r.post.ReplyCount())
}

// UpvoteCount graphql
func (r *postResolver) UpvoteCount() int32 {
	return int32(r.post.UpvoteCount())
}

// DownvoteCount graphql
func (r *postResolver) DownvoteCount() int32 {
	return int32(r.post.DownvoteCount())
}

//Upvoted bool
func (r *postResolver) Upvoted() bool {
	accountID := r.pr.Context.Value("request").(map[string]string)["id"]
	return r.post.IsUpvoted(accountID)
}

//Downvoted bool
func (r *postResolver) Downvoted() bool {
	accountID := r.pr.Context.Value("request").(map[string]string)["id"]
	return r.post.IsDownvoted(accountID)
}

// Parent graphql
func (r *postResolver) Parent() *postResolver {
	return nil
}

//Repost graphql
func (r *postResolver) Repost() *postResolver {
	if r.post.Repost() == nil {
		return nil
	}
	return &postResolver{r.post.Repost(), r.pr}
}

// Child graphql
func (r *postResolver) Child(args struct {
	First  *int32
	After  *graphql.ID
	Before *graphql.ID
	Sort   *int32
}) *postConnectionResolver {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := "000000000000000000000000"
	postList := r.pr.Interactor.PostChild(r.post.ID(), first, after)
	return &postConnectionResolver{postList, r.pr}
}

// MenfessPostConnectionResolver graphql
type postConnectionResolver struct {
	menfessPostList []entity.Post
	pr              *resolver
}

// Edges graphql
func (r *postConnectionResolver) Edges() []*postResolver {
	var menfessPostResolverList []*postResolver
	for _, post := range r.menfessPostList {
		menfessPostResolverList = append(menfessPostResolverList, &postResolver{post, r.pr})
	}
	return menfessPostResolverList
}

// PageInfo graphql
func (r *postConnectionResolver) PageInfo() *PageInfoResolver {
	var nodeList []node
	for _, node := range r.menfessPostList {
		nodeList = append(nodeList, node)
	}
	return &PageInfoResolver{nodeList}
}
