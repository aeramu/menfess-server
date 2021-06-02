package menfess

import (
	"context"

	"github.com/aeramu/menfess-server/internal/feedback/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRepository(db *mongo.Database) service.MenfessRequestRepository {
	return &repo{
		coll: db.Collection("menfessRequest"),
	}
}

type repo struct {
	coll *mongo.Collection
}

func (r *repo) Save(ctx context.Context, m service.MenfessRequest) error {
	menfess := encode(m)
	update := bson.D{{"$set", menfess}}
	filter := bson.D{{"_id", menfess.ID}}
	opt := options.Update().SetUpsert(true)

	if _, err := r.coll.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	}
	return nil
}

type Menfess struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string
	Desc   string
	Marked bool
	UserID primitive.ObjectID `bson:"userID"`
}

func encode(m service.MenfessRequest) *Menfess {
	return &Menfess{
		ID:     objectID(m.ID),
		Name:   m.Name,
		Desc:   m.Desc,
		Marked: m.Marked,
		UserID: objectID(m.UserID),
	}
}

func objectID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}
