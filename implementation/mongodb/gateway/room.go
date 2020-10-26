package gateway

import (
	"github.com/aeramu/menfess-server/post/service"
	"github.com/aeramu/menfess-server/room"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Room extract model
type Room struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string
	Avatar string
}

//Rooms array of room
type Rooms []*Room

//Entity convert room to entity
func (m *Room) Entity() room.Room {
	return service.RoomConstructor{
		ID:     m.ID.Hex(),
		Name:   m.Name,
		Avatar: m.Avatar,
	}.New()
}

//Entity convert array of room to entity
func (rooms Rooms) Entity() []room.Room {
	var entityList []room.Room
	for _, room := range rooms {
		entityList = append(entityList, room.Entity())
	}
	return entityList
}
