package repository

//TODO: project reply count and delete reply count on db model
//TODO: make bson.D function or aggregation stage function

import (
	"context"
	"github.com/aeramu/menfess-server/post/service"

	"github.com/aeramu/menfess-server/implementation/mongodb/gateway"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lookupRoom = d("$lookup", bson.D{
	e("from", "room"),
	e("localField", "roomID"),
	e("foreignField", "_id"),
	e("as", "room"),
})

var lookupRepost = d("$lookup", bson.D{
	e("from", "post"),
	e("localField", "repostID"),
	e("foreignField", "_id"),
	e("as", "repost"),
})

var lookupRepostRoom = d("$lookup", bson.D{
	e("from", "room"),
	e("localField", "repost.roomID"),
	e("foreignField", "_id"),
	e("as", "repostRoom"),
})

func (repo *repo) GetPostByID(hexID string) service.Post {
	id, _ := primitive.ObjectIDFromHex(hexID)

	filter := d("_id", id)

	match := d("$match", filter)
	cursor, _ := repo.post.Aggregate(context.TODO(), mongo.Pipeline{match, lookupRoom, lookupRepost, lookupRepostRoom})

	var posts []*gateway.Post
	cursor.All(context.TODO(), &posts)

	if len(posts) == 0 {
		return nil
	}
	return posts[0].Entity()
}

func (repo *repo) GetPostListByParentID(parentID string, first int, after string, ascSort bool) []service.Post {
	parentid, _ := primitive.ObjectIDFromHex(parentID)
	afterid, _ := primitive.ObjectIDFromHex(after)
	comparator := "$lt"
	sort := -1
	if ascSort {
		comparator = "$gt"
		sort = 1
	}
	filter := d("$and", bson.A{
		d("parentID", parentid),
		d("_id", d(comparator, afterid)),
	})

	sortOpt := d("$sort", d("_id", sort))
	limit := d("$limit", int64(first))
	match := d("$match", filter)
	cursor, _ := repo.post.Aggregate(context.TODO(),
		mongo.Pipeline{sortOpt, match, limit, lookupRoom, lookupRepost, lookupRepostRoom})

	var posts gateway.Posts
	cursor.All(context.TODO(), &posts)
	return posts.Entity()
}

func (repo *repo) GetPostListByRoomIDs(roomIDs []string, first int, after string, ascSort bool) []service.Post {
	roomids := gateway.IDsFromHex(roomIDs)
	afterid, _ := primitive.ObjectIDFromHex(after)
	comparator := "$lt"
	sort := -1
	if ascSort {
		comparator = "$gt"
		sort = 1
	}
	filter := d("$and", bson.A{
		d("roomID", d("$in", roomids)),
		d("_id", d(comparator, afterid)),
	})

	sortOpt := d("$sort", d("_id", sort))
	limit := d("$limit", int64(first))
	match := d("$match", filter)
	cursor, _ := repo.post.Aggregate(context.TODO(),
		mongo.Pipeline{sortOpt, match, limit, lookupRoom, lookupRepost, lookupRepostRoom})

	var posts gateway.Posts
	cursor.All(context.TODO(), &posts)
	return posts.Entity()
}

func (repo *repo) PutPost(name string, avatar string, body string, parentID string, repostID string, roomID string) service.Post {
	post := gateway.NewPostModel(name, avatar, body, parentID, repostID, roomID)
	filter := d("_id", post.ParentID)
	update := d("$inc", d("replyCount", 1))
	option := options.BulkWrite().SetOrdered(false)
	models := []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(post),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update),
	}
	repo.post.BulkWrite(context.TODO(), models, option)

	return repo.GetPostByID(post.ID.Hex())
}

func (repo *repo) UpdateUpvoterIDs(postID string, accountID string, exist bool) {
	operator := "$set"
	if exist {
		operator = "$unset"
	}
	postid, _ := primitive.ObjectIDFromHex(postID)

	filter := d("_id", postid)
	update := d(operator, d("upvoterIDs."+accountID, true))

	repo.post.UpdateOne(context.TODO(), filter, update)
}

func (repo *repo) UpdateDownvoterIDs(postID string, accountID string, exist bool) {
	operator := "$set"
	if exist {
		operator = "$unset"
	}
	postid, _ := primitive.ObjectIDFromHex(postID)

	filter := d("_id", postid)
	update := d(operator, d("downvoterIDs."+accountID, true))

	repo.post.UpdateOne(context.TODO(), filter, update)
}
