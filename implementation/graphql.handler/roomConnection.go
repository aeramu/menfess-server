package resolver

import "github.com/aeramu/menfess-server/entity"

//RoomConnection graphql
type RoomConnection interface {
	Edges() []Room
	PageInfo() PageInfo
}

// MenfessRoomConnectionResolver graphql
type roomConnection struct {
	menfessRoomList []entity.Room
	pr              *resolver
}

// Edges graphql
func (r *roomConnection) Edges() []Room {
	var menfessRoomResolverList []Room
	for _, elem := range r.menfessRoomList {
		menfessRoomResolverList = append(menfessRoomResolverList, &room{elem, r.pr})
	}
	return menfessRoomResolverList
}

// PageInfo graphql
func (r *roomConnection) PageInfo() PageInfo {
	var nodeList []interface{ ID() string }
	for _, node := range r.menfessRoomList {
		nodeList = append(nodeList, node)
	}
	return &pageInfo{nodeList}
}
