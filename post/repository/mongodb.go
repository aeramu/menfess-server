package repository

import (
	"context"
	"fmt"
	"github.com/aeramu/menfess-server/post/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client *mongo.Client

func NewRepository() service.Repository {
	if client == nil {
		client, _ = mongo.Connect(context.Background(), options.Client().ApplyURI(
			"mongodb+srv://admin:admin@qiup-wrbox.mongodb.net/",
		))
	}
	return &repo{
		coll:   client.Database("menfess").Collection("post"),
	}
}

type repo struct{
	coll *mongo.Collection
}

func (r *repo) Save(p service.Post) error {
	post := encode(p)
	update := bson.D{{"$set", post}}
	filter := bson.D{{"_id", post.ID}}
	opt := options.Update().SetUpsert(true)

	if _, err := r.coll.UpdateOne(context.TODO(), filter, update, opt); err != nil{
		return err
	}
	return nil
}

func (r *repo) FindByID(id string) (*service.Post, error) {
	filter := bson.D{{"_id", objectID(id)}}

	cursor, err := r.coll.Find(context.TODO(), filter)
	if err != nil{
		return nil, err
	}

	var posts []Post
	if err := cursor.All(context.TODO(), &posts); err != nil{
		return nil, err
	}

	if len(posts) == 0{
		return nil, nil
	}
	return posts[0].decode(), nil
}

func (r *repo) FindByParentID(id string, first int, after string, sort bool) (*[]service.Post, error) {
	comparator := "$gt"
	sortOpt := bson.D{{"_id", 1}}
	if sort {
		comparator = "$lt"
		sortOpt = bson.D{{"_id", -1}}
	}

	filter := bson.D{{"$and", bson.A{
		bson.D{{"parentID", objectID(id)}},
		bson.D{{"_id",bson.D{{comparator, objectID(after)}}}},
	}}}
	opt := options.Find().SetLimit(int64(first)).SetSort(sortOpt)

	cursor, err := r.coll.Find(context.TODO(), filter, opt)
	if err != nil{
		return nil, err
	}

	var posts Posts
	if err := cursor.All(context.TODO(), &posts); err != nil{
		return nil, err
	}
	fmt.Println()

	return posts.decode(), nil
}

func (r *repo) FindByRoomID(id string, first int, after string, sort bool) (*[]service.Post, error) {
	comparator := "$gt"
	sortOpt := bson.D{{"_id", 1}}
	if sort {
		comparator = "$lt"
		sortOpt = bson.D{{"_id", -1}}
	}

	filter := bson.D{{"$and", bson.A{
		bson.D{{"parentID", objectID("")}},
		bson.D{{"roomID", objectID(id)}},
		bson.D{{"_id",bson.D{{comparator, objectID(after)}}}},
	}}}
	opt := options.Find().SetLimit(int64(first)).SetSort(sortOpt)

	cursor, err := r.coll.Find(context.TODO(), filter, opt)
	if err != nil{
		return nil, err
	}

	var posts Posts
	if err := cursor.All(context.TODO(), &posts); err != nil{
		return nil, err
	}

	return posts.decode(), nil
}

type Post struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	Avatar       string
	Body         string
	ParentID     primitive.ObjectID `bson:"parentID"`
	RepostID     primitive.ObjectID `bson:"repostID"`
	RoomID       primitive.ObjectID `bson:"roomID"`
	UpvoterIDs   map[string]bool    `bson:"upvoterIDs"`
	DownvoterIDs map[string]bool    `bson:"downvoterIDs"`
	ReplyCount   int                `bson:"replyCount"`
}

func (p Post) decode() *service.Post{
	return &service.Post{
		ID:           p.ID.Hex(),
		Timestamp:    int(p.ID.Timestamp().Unix()),
		Name:         p.Name,
		Avatar:       p.Avatar,
		Body:         p.Body,
		ParentID:     p.ParentID.Hex(),
		RepostID:     p.RepostID.Hex(),
		RoomID:       p.RoomID.Hex(),
		UpvoterIDs:   p.UpvoterIDs,
		DownvoterIDs: p.DownvoterIDs,
		ReplyCount:   p.ReplyCount,
	}
}

type Posts []Post

func (p Posts) decode() *[]service.Post{
	var posts []service.Post
	for _, post := range p {
		posts = append(posts, *post.decode())
	}
	return &posts
}

func encode(p service.Post) *Post{
	return &Post{
		ID:           objectID(p.ID),
		Name:         p.Name,
		Avatar:       p.Avatar,
		Body:         p.Body,
		ParentID:     objectID(p.ParentID),
		RepostID:     objectID(p.RepostID),
		RoomID:       objectID(p.RoomID),
		UpvoterIDs:   p.UpvoterIDs,
		DownvoterIDs: p.DownvoterIDs,
		ReplyCount:   p.ReplyCount,
	}
}

func objectID(hex string) primitive.ObjectID{
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Println("Parse Mongo ID Error", err)
	}
	return id
}
