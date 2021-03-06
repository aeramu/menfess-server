package repository

import (
	"context"
	"github.com/aeramu/menfess-server/room"

	"github.com/aeramu/menfess-server/implementation/mongodb/gateway"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *repo) GetRoomList() []room.Room {
	filter := bson.D{{}}
	cursor, _ := repo.room.Find(context.TODO(), filter)

	var rooms gateway.Rooms
	cursor.All(context.TODO(), &rooms)
	return rooms.Entity()
}

func (repo *repo) GetRoom(id string) room.Room {
	objectID, _ := primitive.ObjectIDFromHex(id)

	filter := d("_id", objectID)
	var room gateway.Room
	repo.room.FindOne(context.TODO(), filter).Decode(&room)

	if room.ID.IsZero() {
		return nil
	}
	return room.Entity()
}
