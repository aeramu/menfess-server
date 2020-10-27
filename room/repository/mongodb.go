package repository

import (
	"context"
	"fmt"
	"github.com/aeramu/menfess-server/room/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client *mongo.Client

func NewRepository() service.Repository {
	var err error
	if client == nil {
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(
			"mongodb+srv://admin:admin@qiup-wrbox.mongodb.net/",
		))
	}
	if err != nil{
		log.Println("DB Connect Error:", err)
		return nil
	}
	return &repo{
		coll:   client.Database("menfess").Collection("room"),
	}
}

type repo struct{
	coll *mongo.Collection
}

func (r *repo) Save(rm service.Room) error {
	room := encode(rm)
	update := bson.D{{"$set", room}}
	filter := bson.D{{"_id", room.ID}}
	opt := options.Update().SetUpsert(true)

	if _, err := r.coll.UpdateOne(context.TODO(), filter, update, opt); err != nil{
		return err
	}
	return nil
}

func (r *repo) FindByID(id string) (*service.Room, error) {
	filter := bson.D{{"_id", objectID(id)}}

	cursor, err := r.coll.Find(context.TODO(), filter)
	if err != nil{
		return nil, err
	}

	var posts []Room
	if err := cursor.All(context.TODO(), &posts); err != nil{
		return nil, err
	}

	if len(posts) == 0{
		return nil, nil
	}
	return posts[0].decode(), nil
}

func (r *repo) FindByStatus(status bool) (*[]service.Room, error) {
	filter := bson.D{{"status", status}}

	cursor, err := r.coll.Find(context.TODO(), filter)
	if err != nil{
		return nil, err
	}

	var rooms Rooms
	if err := cursor.All(context.TODO(), &rooms); err != nil{
		return nil, err
	}
	fmt.Println()

	return rooms.decode(), nil
}

func (r *repo) FindAll() (*[]service.Room, error) {
	filter := bson.D{{}}

	cursor, err := r.coll.Find(context.TODO(), filter)
	if err != nil{
		return nil, err
	}

	var rooms Rooms
	if err := cursor.All(context.TODO(), &rooms); err != nil{
		return nil, err
	}

	return rooms.decode(), nil
}

type Room struct {
	ID primitive.ObjectID `bson:"_id"`
	Name string
	Desc string
	Avatar string
	Status bool
}

func (r Room) decode() *service.Room{
	return &service.Room{
		ID:           r.ID.Hex(),
		Name:         r.Name,
		Desc: r.Desc,
		Avatar: r.Avatar,
		Status: r.Status,
	}
}

type Rooms []Room

func (p Rooms) decode() *[]service.Room{
	var rooms []service.Room
	for _, room := range p {
		rooms = append(rooms, *room.decode())
	}
	return &rooms
}

func encode(r service.Room) *Room {
	return &Room{
		ID:           objectID(r.ID),
		Name:         r.Name,
		Desc: r.Desc,
		Avatar: r.Avatar,
		Status: r.Status,
	}
}

func objectID(hex string) primitive.ObjectID{
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}
