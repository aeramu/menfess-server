package repository

import (
	"context"

	"github.com/aeramu/menfess-server/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client = nil

//New MenfessPostRepo Constructor
func New() usecase.Repository {
	if client == nil {
		client, _ = mongo.Connect(context.Background(), options.Client().ApplyURI(
			"mongodb+srv://admin:admin@qiup-wrbox.mongodb.net/",
		))
	}
	return &repo{
		client: client,
		db:     client.Database("menfess"),
		post:   client.Database("menfess").Collection("post"),
		room:   client.Database("menfess").Collection("room"),
	}
}

//Disconnect disconnect
func Disconnect() {
	client.Disconnect(context.TODO())
}

type repo struct {
	client *mongo.Client
	db     *mongo.Database
	post   *mongo.Collection
	room   *mongo.Collection
}
