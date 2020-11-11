package repository

import (
	"context"
	"log"

	"github.com/aeramu/menfess-server/user/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func NewRepository() service.Repository {
	var err error
	if client == nil {
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(
			"mongodb+srv://admin:admin@qiup-wrbox.mongodb.net/",
		))
	}
	if err != nil {
		log.Println("DB Connect Error:", err)
		return nil
	}
	return &repo{
		coll: client.Database("menfessv2").Collection("user"),
	}
}

type repo struct {
	coll *mongo.Collection
}

func (r *repo) Save(p service.User) error {
	user := encode(p)
	update := bson.D{{"$set", user}}
	filter := bson.D{{"_id", user.ID}}
	opt := options.Update().SetUpsert(true)

	if _, err := r.coll.UpdateOne(context.TODO(), filter, update, opt); err != nil {
		return err
	}
	return nil
}

func (r *repo) FindByID(id string) (*service.User, error) {
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

func (r *repo) FindByEmail(email string) (*service.User, error) {
	filter := bson.D{{"email", email}}

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

func (r *repo) FindByType(t string) (*[]service.User, error) {
	filter := bson.D{{"type", t}}

	sort := bson.D{{"name", 1}}
	opt := options.Find().SetSort(sort)

	cursor, err := r.coll.Find(context.TODO(), filter, opt)
	if err != nil {
		return nil, err
	}

	var users Users
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users.decode(), nil
}

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string
	Password  string
	Name      string
	Avatar    string
	Bio       string
	PushToken map[string]bool `bson:"pushToken"`
}

func (u User) decode() *service.User {
	return &service.User{
		ID:        u.ID.Hex(),
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		PushToken: u.PushToken,
	}
}

type Users []User

func (p Users) decode() *[]service.User {
	var users []service.User
	for _, user := range p {
		users = append(users, *user.decode())
	}
	return &users
}

func encode(u service.User) *User {
	return &User{
		ID:        objectID(u.ID),
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		PushToken: u.PushToken,
	}
}

func objectID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}
