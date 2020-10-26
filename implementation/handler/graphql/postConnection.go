package resolver

import (
	resolver2 "github.com/aeramu/menfess-server/gateway/resolver"
	"github.com/aeramu/menfess-server/post/service"
)

//PostConnection graphql
type PostConnection interface {
	Edges() []resolver2.Post
	PageInfo() PageInfo
}

type postConnection struct {
	menfessPostList []service.Post
	pr              *resolver
}

// Edges graphql
func (r *postConnection) Edges() []resolver2.Post {
	var menfessPostResolverList []resolver2.Post
	for _, elem := range r.menfessPostList {
		menfessPostResolverList = append(menfessPostResolverList, &resolver2.post{elem, r.pr})
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
