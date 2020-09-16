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
	UpvoterIDs   map[string]bool `bson:"upvoterIDs"`
	DownvoterIDs map[string]bool `bson:"downvoterIDs"`
	ReplyCount   int             `bson:"replyCount"`
	Repost       []Post          `bson:"repost"`
	Room         []Room          `bson:"room"`
}

type Posts []*Post

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
		Repost:       m.Repost[0].Entity(),
		Room:         m.Room[0].Entity(),
	}.New()
}

func (posts Posts) Entity() []entity.Post {
	var entityList []entity.Post
	for _, post := range posts {
		entityList = append(entityList, post.Entity())
	}
	return entityList
}
