package resolver

import "github.com/aeramu/menfess-server/entity"

//PostConnection graphql
type PostConnection interface {
	Edges() []Post
	PageInfo() PageInfo
}

// MenfessPostConnectionResolver graphql
type postConnectionResolver struct {
	menfessPostList []entity.Post
	pr              *resolver
}

// Edges graphql
func (r *postConnectionResolver) Edges() []Post {
	var menfessPostResolverList []Post
	for _, elem := range r.menfessPostList {
		menfessPostResolverList = append(menfessPostResolverList, &post{elem, r.pr})
	}
	return menfessPostResolverList
}

// PageInfo graphql
func (r *postConnectionResolver) PageInfo() PageInfo {
	var nodeList []interface{ ID() string }
	for _, node := range r.menfessPostList {
		nodeList = append(nodeList, node)
	}
	return pageInfo{nodeList}
}
