package feedback

import (
	"context"

	"github.com/aeramu/menfess-server/internal/feedback/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRepository(db *mongo.Database) service.FeedbackRepository {
	return &repo{
		coll: db.Collection("feedback"),
	}
}

type repo struct {
	coll *mongo.Collection
}

func (r *repo) Save(ctx context.Context, f service.Feedback) error {
	feedback := encode(f)
	update := bson.D{{"$set", feedback}}
	filter := bson.D{{"_id", feedback.ID}}
	opt := options.Update().SetUpsert(true)

	if _, err := r.coll.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	}
	return nil
}

type Feedback struct {
	ID       primitive.ObjectID `bson:"_id"`
	Feedback string
	Marked   bool
	UserID   primitive.ObjectID `bson:"userID"`
}

func encode(f service.Feedback) *Feedback {
	return &Feedback{
		ID:       objectID(f.ID),
		Feedback: f.Feedback,
		Marked:   f.Marked,
		UserID:   objectID(f.UserID),
	}
}

func objectID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}
