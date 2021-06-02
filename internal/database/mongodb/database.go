package mongodb

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dblock sync.Once
	client *mongo.Client
)

func NewDatabase() *mongo.Database {
	var err error
	dblock.Do(func() {
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(
			"mongodb+srv://admin:admin@qiup-wrbox.mongodb.net/",
		))
	})

	if err != nil {
		log.Println("Failed Connect MongoDB:", err)
	}

	return client.Database("menfessv2")
}
