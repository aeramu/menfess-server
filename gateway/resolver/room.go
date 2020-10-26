package resolver

import "github.com/graph-gophers/graphql-go"

type Room struct {}
func (r Room) ID() graphql.ID{
	return ""
}
func (r Room) Name() string{
	return ""
}
func (r Room) Avatar() string {
	return ""
}

type RoomConnection struct{
	rooms []*Room
	root  *resolver
}
func (r *RoomConnection) Edges() []*Room {
	return r.rooms
}
func (r *RoomConnection) PageInfo() PageInfo {
	var nodeList PageInfo
	for _, node := range r.rooms {
		nodeList = append(nodeList, node)
	}
	return nodeList
}
