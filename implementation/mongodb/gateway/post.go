package gateway

import (
	"github.com/aeramu/menfess-server/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
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

type Posts []*Post

func NewPost(name string, avatar string, body string, parentID string, repostID string, roomID string) *Post {
	parentid, _ := primitive.ObjectIDFromHex(parentID)
	repostid, _ := primitive.ObjectIDFromHex(repostID)
	roomid, _ := primitive.ObjectIDFromHex(roomID)
	return &Post{
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

func PostFromEntity(e entity.Post) *Post {
	id, _ := primitive.ObjectIDFromHex(e.ID())
	parentID, _ := primitive.ObjectIDFromHex(e.ParentID())
	repostID, _ := primitive.ObjectIDFromHex(e.RepostID())
	roomID, _ := primitive.ObjectIDFromHex(e.RoomID())
	return &Post{
		ID:           id,
		Name:         e.Name(),
		Avatar:       e.Avatar(),
		Body:         e.Body(),
		UpvoterIDs:   e.UpvoterIDs(),
		DownvoterIDs: e.DownvoterIDs(),
		ReplyCount:   e.ReplyCount(),
		ParentID:     parentID,
		RepostID:     repostID,
		RoomID:       roomID,
	}
}

func (m *Post) Entity() entity.Post {
	return entity.PostConstructor{
		ID:           m.ID.Hex(),
		Timestamp:    int(m.ID.Timestamp().Unix()),
		Name:         m.Name,
		Avatar:       m.Avatar,
		Body:         m.Body,
		UpvoterIDs:   m.UpvoterIDs,
		DownvoterIDs: m.DownvoterIDs,
		ReplyCount:   m.ReplyCount,
		ParentID:     m.ParentID.Hex(),
		RepostID:     m.RepostID.Hex(),
		RoomID:       m.RoomID.Hex(),
	}.New()
}

func (posts Posts) Entity() []entity.Post {
	var entityList []entity.Post
	for _, post := range posts {
		entityList = append(entityList, post.Entity())
	}
	return entityList
}
