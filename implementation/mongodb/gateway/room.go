package gateway

import (
	"github.com/aeramu/menfess-server/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Room db model
type Room struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string
	Avatar string
}

//Rooms arrat of room's db model
type Rooms []*Room

//Entity convert room db model to entity
func (m *Room) Entity() entity.Room {
	return entity.RoomConstructor{
		ID:     m.ID.Hex(),
		Name:   m.Name,
		Avatar: m.Avatar,
	}.New()
}

//Entity convert array of room db model to entity
func (rooms Rooms) Entity() []entity.Room {
	var entityList []entity.Room
	for _, room := range rooms {
		entityList = append(entityList, room.Entity())
	}
	return entityList
}
