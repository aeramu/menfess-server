package gateway

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//PostModel post with db format
type PostModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	Avatar       string
	Body         string
	UpvoterIDs   map[string]bool    `bson:"upvoterIDs"`
	DownvoterIDs map[string]bool    `bson:"downvoterIDs"`
	ReplyCount   int                `bson:"replyCount"`
	ParentID     primitive.ObjectID `bson:"parentID"`
	RepostID     primitive.ObjectID `bson:"repostID"`
	RoomID       primitive.ObjectID `bson:"roomID"`
}

//NewPostModel create new post with db format
func NewPostModel(name string, avatar string, body string, parentID string, repostID string, roomID string) *PostModel {
	parentid, _ := primitive.ObjectIDFromHex(parentID)
	repostid, _ := primitive.ObjectIDFromHex(repostID)
	roomid, _ := primitive.ObjectIDFromHex(roomID)
	return &PostModel{
		ID:           primitive.NewObjectID(),
		Name:         name,
		Avatar:       avatar,
		Body:         body,
		UpvoterIDs:   map[string]bool{},
		DownvoterIDs: map[string]bool{},
		ReplyCount:   0,
		ParentID:     parentid,
		RepostID:     repostid,
		RoomID:       roomid,
	}
}
