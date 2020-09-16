package repository

import (
	"context"

	"github.com/aeramu/menfess-server/entity"
	"github.com/aeramu/menfess-server/implementation/mongodb/gateway"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *repo) GetPostByID(hexID string) entity.Post {
	id, _ := primitive.ObjectIDFromHex(hexID)

	filter := bson.D{{"_id", id}}
	var post gateway.Post
	repo.post.FindOne(context.TODO(), filter).Decode(&post)

	if post.ID.IsZero() {
		return nil
	}
	return post.Entity()
}

func (repo *repo) GetPostListByParentID(parentID string, first int, after string, ascSort bool) []entity.Post {
	parentid, _ := primitive.ObjectIDFromHex(parentID)
	afterid, _ := primitive.ObjectIDFromHex(after)
	comparator := "$lt"
	sort := -1
	if ascSort {
		comparator = "$gt"
		sort = 1
	}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"parentID", parentid}},
			bson.D{
				{"_id", bson.D{
					{comparator, afterid},
				}},
			},
		}},
	}
	sortOpt := bson.D{{"_id", sort}}
	option := options.Find().SetLimit(int64(first)).SetSort(sortOpt)
	//TODO: lookup room dan repost dan parent?

	// lookup := bson.D{
	// 	{"$lookup", bson.D{
	// 		{"from", "room"},
	// 		{"localField", "roomID"},
	// 		{"foreignField", "_id"},
	// 		{"as", "room"},
	// 	}},
	// }

	// cursor2, _ := repo.collection.Aggregate(context.TODO(), mongo.Pipeline{lookup})
	cursor, _ := repo.post.Find(context.TODO(), filter, option)

	var posts gateway.Posts
	cursor.All(context.TODO(), &posts)
	return posts.Entity()
}

func (repo *repo) GetPostListByRoomIDs(roomIDs []string, first int, after string, ascSort bool) []entity.Post {
	roomids := gateway.IDsFromHex(roomIDs)
	afterid, _ := primitive.ObjectIDFromHex(after)
	comparator := "$lt"
	sort := -1
	if ascSort {
		comparator = "$gt"
		sort = 1
	}

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"roomID", bson.D{
				{"$in", roomids},
			}}},
			bson.D{{"_id", bson.D{
				{comparator, afterid},
			}}},
		}},
	}
	sortOpt := bson.D{{"_id", sort}}
	option := options.Find().SetLimit(int64(first)).SetSort(sortOpt)
	cursor, _ := repo.post.Find(context.TODO(), filter, option)

	var posts gateway.Posts
	cursor.All(context.TODO(), &posts)
	return posts.Entity()
}

func (repo *repo) PutPost(name string, avatar string, body string, parentID string, repostID string, roomID string) entity.Post {
	post := gateway.NewPost(name, avatar, body, parentID, repostID, roomID)
	filter := bson.D{{"_id", post.ParentID}}
	update := bson.D{
		{"$inc", bson.D{
			{"replyCount", 1},
		}},
	}
	option := options.BulkWrite().SetOrdered(false)
	models := []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(post),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true),
	}
	repo.post.BulkWrite(context.TODO(), models, option)

	return post.Entity()
}

func (repo *repo) UpdateUpvoterIDs(postID string, accountID string, exist bool) {
	operator := "$set"
	if exist {
		operator = "$unset"
	}
	postid, _ := primitive.ObjectIDFromHex(postID)

	filter := bson.D{{"_id", postid}}
	update := bson.D{
		{operator, bson.D{
			{"upvoterIDs." + accountID, true},
		}},
	}
	repo.post.UpdateOne(context.TODO(), filter, update)
}

func (repo *repo) UpdateDownvoterIDs(postID string, accountID string, exist bool) {
	operator := "$set"
	if exist {
		operator = "$unset"
	}
	postid, _ := primitive.ObjectIDFromHex(postID)

	filter := bson.D{{"_id", postid}}
	update := bson.D{
		{operator, bson.D{
			{"downvoterIDs." + accountID, true},
		}},
	}
	repo.post.UpdateOne(context.TODO(), filter, update)
}
