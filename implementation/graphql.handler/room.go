package resolver

import (
	"github.com/aeramu/menfess-server/entity"
	"github.com/graph-gophers/graphql-go"
)

//MenfessRoomResolver graphql
type MenfessRoomResolver struct {
	room entity.Room
	pr   *resolver
}

//ID get
func (r *MenfessRoomResolver) ID() graphql.ID {
	return graphql.ID(r.room.ID())
}

//Name get
func (r *MenfessRoomResolver) Name() string {
	return r.room.Name()
}

// Avatar graphql
func (r *MenfessRoomResolver) Avatar() string {
	return r.room.Avatar()
}

// MenfessRoomConnectionResolver graphql
type MenfessRoomConnectionResolver struct {
	menfessRoomList []entity.Room
	pr              *resolver
}

// Edges graphql
func (r *MenfessRoomConnectionResolver) Edges() []*MenfessRoomResolver {
	var menfessRoomResolverList []*MenfessRoomResolver
	for _, room := range r.menfessRoomList {
		menfessRoomResolverList = append(menfessRoomResolverList, &MenfessRoomResolver{room, r.pr})
	}
	return menfessRoomResolverList
}

// PageInfo graphql
func (r *MenfessRoomConnectionResolver) PageInfo() PageInfo {
	var nodeList []interface{ ID() string }
	for _, node := range r.menfessRoomList {
		nodeList = append(nodeList, node)
	}
	return &pageInfo{nodeList}
}
