package gateway

import (
	"github.com/aeramu/menfess-server/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//PostModel db model
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

//NewPostModel create new post db model
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

//Post extract model
type Post struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	Avatar       string
	Body         string
	UpvoterIDs   map[string]bool `bson:"upvoterIDs"`
	DownvoterIDs map[string]bool `bson:"downvoterIDs"`
	ReplyCount   int             `bson:"replyCount"`
	Repost       []Post          `bson:"repost"`
	Room         []Room          `bson:"room"`
	RepostRoom   []Room          `bson:"repostRoom"`
}

//Posts list of post
type Posts []*Post

//Entity convert post to entity
func (m *Post) Entity() entity.Post {
	var repost entity.Post = nil
	if len(m.Repost) > 0 {
		m.Repost[0].Room = m.RepostRoom
		repost = m.Repost[0].Entity()
	}
	var room entity.Room = nil
	if len(m.Room) > 0 {
		room = m.Room[0].Entity()
	}

	return entity.PostConstructor{
		ID:           m.ID.Hex(),
		Timestamp:    int(m.ID.Timestamp().Unix()),
		Name:         m.Name,
		Avatar:       m.Avatar,
		Body:         m.Body,
		UpvoterIDs:   m.UpvoterIDs,
		DownvoterIDs: m.DownvoterIDs,
		ReplyCount:   m.ReplyCount,
		Repost:       repost,
		Room:         room,
	}.New()
}

//Entity convert array of post to entity
func (posts Posts) Entity() []entity.Post {
	var entityList []entity.Post
	for _, post := range posts {
		entityList = append(entityList, post.Entity())
	}
	return entityList
}
