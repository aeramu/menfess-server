package repository

import (
	"context"

	"github.com/aeramu/menfess-server/entity"
	"github.com/aeramu/menfess-server/implementation/mongodb/gateway"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *menfessRepo) GetRoomList() []entity.Room {
	filter := bson.D{{}}
	cursor, _ := repo.client.Database("menfess").Collection("room").Find(context.TODO(), filter)

	var rooms gateway.Rooms
	cursor.All(context.TODO(), &rooms)
	return rooms.Entity()
}

func (repo *menfessRepo) GetRoom(id string) entity.Room {
	objectID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{"_id", objectID}}
	var room gateway.Room
	repo.client.Database("menfess").Collection("room").FindOne(context.TODO(), filter).Decode(&room)

	if room.ID.IsZero() {
		return nil
	}
	return room.Entity()
}
