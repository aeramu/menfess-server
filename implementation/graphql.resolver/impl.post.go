package resolver

import (
	"github.com/aeramu/menfess-server/entity"
	"github.com/graph-gophers/graphql-go"
)

// MenfessPostResolver graphql
type MenfessPostResolver struct {
	post entity.Post
	pr   *Resolver
}

// ID graphql
func (r *MenfessPostResolver) ID() graphql.ID {
	return graphql.ID(r.post.ID())
}

// Timestamp graphql
func (r *MenfessPostResolver) Timestamp() int32 {
	return int32(r.post.Timestamp())
}

// Name graphql
func (r *MenfessPostResolver) Name() string {
	return r.post.Name()
}

// Avatar graphql
func (r *MenfessPostResolver) Avatar() string {
	return r.post.Avatar()
}

// Body graphql
func (r *MenfessPostResolver) Body() string {
	return r.post.Body()
}

//Room graphql
func (r *MenfessPostResolver) Room() string {
	if r.post.Room() == nil {
		return "General"
	}
	return r.post.Room().Name()
}

// ReplyCount graphql
func (r *MenfessPostResolver) ReplyCount() int32 {
	return int32(r.post.ReplyCount())
}

// UpvoteCount graphql
func (r *MenfessPostResolver) UpvoteCount() int32 {
	return int32(r.post.UpvoteCount())
}

// DownvoteCount graphql
func (r *MenfessPostResolver) DownvoteCount() int32 {
	return int32(r.post.DownvoteCount())
}

//Upvoted bool
func (r *MenfessPostResolver) Upvoted() bool {
	accountID := r.pr.Context.Value("request").(map[string]string)["id"]
	return r.post.IsUpvoted(accountID)
}

//Downvoted bool
func (r *MenfessPostResolver) Downvoted() bool {
	accountID := r.pr.Context.Value("request").(map[string]string)["id"]
	return r.post.IsDownvoted(accountID)
}

// Parent graphql
func (r *MenfessPostResolver) Parent() *MenfessPostResolver {
	return nil
}

//Repost graphql
func (r *MenfessPostResolver) Repost() *MenfessPostResolver {
	if r.post.Repost() == nil {
		return nil
	}
	return &MenfessPostResolver{r.post.Repost(), r.pr}
}

// Child graphql
func (r *MenfessPostResolver) Child(args struct {
	First  *int32
	After  *graphql.ID
	Before *graphql.ID
	Sort   *int32
}) *MenfessPostConnectionResolver {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := "000000000000000000000000"
	postList := r.pr.Interactor.PostChild(r.post.ID(), first, after)
	return &MenfessPostConnectionResolver{postList, r.pr}
}

// MenfessPostConnectionResolver graphql
type MenfessPostConnectionResolver struct {
	menfessPostList []entity.Post
	pr              *Resolver
}

// Edges graphql
func (r *MenfessPostConnectionResolver) Edges() []*MenfessPostResolver {
	var menfessPostResolverList []*MenfessPostResolver
	for _, post := range r.menfessPostList {
		menfessPostResolverList = append(menfessPostResolverList, &MenfessPostResolver{post, r.pr})
	}
	return menfessPostResolverList
}

// PageInfo graphql
func (r *MenfessPostConnectionResolver) PageInfo() *PageInfoResolver {
	var nodeList []node
	for _, node := range r.menfessPostList {
		nodeList = append(nodeList, node)
	}
	return &PageInfoResolver{nodeList}
}
