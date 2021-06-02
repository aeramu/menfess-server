package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aeramu/menfess-server/internal/notification/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRepository(db *mongo.Database) service.Repository {
	return &repo{
		coll: db.Collection("notif"),
	}
}

type repo struct {
	coll *mongo.Collection
}

func (r *repo) FindByID(ctx context.Context, id string) (*service.User, error) {
	filter := bson.D{{"_id", objectID(id)}}

	cursor, err := r.coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}
	return users[0].decode(), nil
}

func (r *repo) Save(ctx context.Context, u service.User) error {
	user := encode(u)
	update := bson.D{{"$set", user}}
	filter := bson.D{{"_id", user.ID}}
	opt := options.Update().SetUpsert(true)

	if _, err := r.coll.UpdateOne(context.TODO(), filter, update, opt); err != nil {
		return err
	}
	return nil
}

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	PushToken map[string]bool    `bson:"pushToken"`
}

func (u User) decode() *service.User {
	return &service.User{
		ID:        u.ID.Hex(),
		PushToken: u.PushToken,
	}
}

type Users []User

func (p Users) decode() []service.User {
	var users []service.User
	for _, user := range p {
		users = append(users, *user.decode())
	}
	return users
}

func encode(u service.User) *User {
	return &User{
		ID:        objectID(u.ID),
		PushToken: u.PushToken,
	}
}

func objectID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}
