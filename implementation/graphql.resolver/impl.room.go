package resolver

import (
	"github.com/aeramu/menfess-server/entity"
	"github.com/graph-gophers/graphql-go"
)

//MenfessRoomList graphql
func (r *Resolver) MenfessRoomList() *MenfessRoomConnectionResolver {
	roomList := r.Interactor.RoomList()
	return &MenfessRoomConnectionResolver{roomList, r}
}

//MenfessRoomResolver graphql
type MenfessRoomResolver struct {
	room entity.Room
	pr   *Resolver
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
	pr              *Resolver
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
func (r *MenfessRoomConnectionResolver) PageInfo() *PageInfoResolver {
	var nodeList []node
	for _, node := range r.menfessRoomList {
		nodeList = append(nodeList, node)
	}
	return &PageInfoResolver{nodeList}
}
