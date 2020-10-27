package resolver

import (
	room "github.com/aeramu/menfess-server/room/service"
	"github.com/graph-gophers/graphql-go"
)

type Room struct {
	room.Room
	root *resolver
}
func (r Room) ID() graphql.ID{
	return graphql.ID(r.Room.ID)
}
func (r Room) Name() string{
	return r.Room.Name
}
func (r Room) Desc() string{
	return r.Room.Name
}
func (r Room) Avatar() string {
	return r.Room.Avatar
}
func (r Room) Status() bool{
	return r.Room.Status
}

type RoomConnection struct{
	rooms []Room
	root  *resolver
}
func (r RoomConnection) Edges() []Room {
	return r.rooms
}
func (r RoomConnection) PageInfo() PageInfo {
	var nodeList PageInfo
	for _, node := range r.rooms {
		nodeList = append(nodeList, node)
	}
	return nodeList
}
