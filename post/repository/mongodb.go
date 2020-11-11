package repository

import (
	"context"
	"log"

	"github.com/aeramu/menfess-server/post/service"
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
		coll: client.Database("menfessv2").Collection("post"),
	}
}

type repo struct {
	coll *mongo.Collection
}

func (r *repo) Save(p service.Post) error {
	post := encode(p)
	update := bson.D{{"$set", post}}
	filter := bson.D{{"_id", post.ID}}
	opt := options.Update().SetUpsert(true)

	if _, err := r.coll.UpdateOne(context.TODO(), filter, update, opt); err != nil {
		return err
	}
	return nil
}

func (r *repo) FindByID(id string) (*service.Post, error) {
	filter := bson.D{{"_id", objectID(id)}}

	cursor, err := r.coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var posts []Post
	if err := cursor.All(context.TODO(), &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
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
		bson.D{{"_id", bson.D{{comparator, objectID(after)}}}},
		bson.D{{"reportsCount", bson.D{{"$lt", 2}}}},
	}}}
	opt := options.Find().SetLimit(int64(first)).SetSort(sortOpt)

	cursor, err := r.coll.Find(context.TODO(), filter, opt)
	if err != nil {
		return nil, err
	}

	var posts Posts
	if err := cursor.All(context.TODO(), &posts); err != nil {
		return nil, err
	}

	return posts.decode(), nil
}

func (r *repo) FindByAuthorID(id string, first int, after string, sort bool) (*[]service.Post, error) {
	comparator := "$gt"
	sortOpt := bson.D{{"_id", 1}}
	if sort {
		comparator = "$lt"
		sortOpt = bson.D{{"_id", -1}}
	}

	filter := bson.D{{"$and", bson.A{
		bson.D{{"authorID", objectID(id)}},
		bson.D{{"_id", bson.D{{comparator, objectID(after)}}}},
		bson.D{{"reportsCount", bson.D{{"$lt", 2}}}},
	}}}
	opt := options.Find().SetLimit(int64(first)).SetSort(sortOpt)

	cursor, err := r.coll.Find(context.TODO(), filter, opt)
	if err != nil {
		return nil, err
	}

	var posts Posts
	if err := cursor.All(context.TODO(), &posts); err != nil {
		return nil, err
	}

	return posts.decode(), nil
}

type Post struct {
	ID           primitive.ObjectID `bson:"_id"`
	Body         string
	AuthorID     primitive.ObjectID `bson:"authorID"`
	UserID       primitive.ObjectID `bson:"userID"`
	ParentID     primitive.ObjectID `bson:"parentID"`
	LikeIDs      map[string]bool    `bson:"likeIDs"`
	RepliesCount int                `bson:"repliesCount"`
	ReportsCount int                `bson:"reportsCount"`
}

func (p Post) decode() *service.Post {
	parentID := p.ParentID.Hex()
	if p.ParentID.IsZero() {
		parentID = ""
	}
	userID := p.UserID.Hex()
	if p.UserID.IsZero() {
		userID = ""
	}
	return &service.Post{
		ID:           p.ID.Hex(),
		Timestamp:    int(p.ID.Timestamp().Unix()),
		Body:         p.Body,
		AuthorID:     p.AuthorID.Hex(),
		UserID:       userID,
		ParentID:     parentID,
		LikeIDs:      p.LikeIDs,
		RepliesCount: p.RepliesCount,
		ReportsCount: p.ReportsCount,
	}
}

type Posts []Post

func (p Posts) decode() *[]service.Post {
	var posts []service.Post
	for _, post := range p {
		posts = append(posts, *post.decode())
	}
	return &posts
}

func encode(p service.Post) *Post {
	return &Post{
		ID:           objectID(p.ID),
		Body:         p.Body,
		AuthorID:     objectID(p.AuthorID),
		UserID:       objectID(p.UserID),
		ParentID:     objectID(p.ParentID),
		LikeIDs:      p.LikeIDs,
		RepliesCount: p.RepliesCount,
		ReportsCount: p.ReportsCount,
	}
}

func objectID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}
