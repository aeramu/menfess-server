package gateway

import (
	"github.com/aeramu/menfess-server/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string
}

type Rooms []*Room

func (m *Room) Entity() entity.Room {
	return entity.RoomConstructor{
		ID:   m.ID.Hex(),
		Name: m.Name,
	}.New()
}

func (rooms Rooms) Entity() []entity.Room {
	var entityList []entity.Room
	for _, room := range rooms {
		entityList = append(entityList, room.Entity())
	}
	return entityList
}
