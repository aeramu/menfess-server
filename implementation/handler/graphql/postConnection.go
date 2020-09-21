package resolver

import "github.com/aeramu/menfess-server/entity"

//PostConnection graphql
type PostConnection interface {
	Edges() []Post
	PageInfo() PageInfo
}

type postConnection struct {
	menfessPostList []entity.Post
	pr              *resolver
}

// Edges graphql
func (r *postConnection) Edges() []Post {
	var menfessPostResolverList []Post
	for _, elem := range r.menfessPostList {
		menfessPostResolverList = append(menfessPostResolverList, &post{elem, r.pr})
	}
	return menfessPostResolverList
}

// PageInfo graphql
func (r *postConnection) PageInfo() PageInfo {
	var nodeList []interface{ ID() string }
	for _, node := range r.menfessPostList {
		nodeList = append(nodeList, node)
	}
	return pageInfo(nodeList)
}
