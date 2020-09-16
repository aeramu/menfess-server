package usecase

import (
	"github.com/aeramu/menfess-server/entity"
)

type interactor struct {
	repo Repository
}

//Repository interface
type Repository interface {
	GetPostByID(id string) entity.Post
	GetPostListByParentID(parentID string, first int, after string, ascSort bool) []entity.Post
	GetPostListByRoomIDs(roomIDs []string, first int, after string, ascSort bool) []entity.Post
	PutPost(name string, avatar string, body string, parentID string, repostID string, roomID string) entity.Post
	UpdateUpvoterIDs(postID string, accountID string, exist bool)
	UpdateDownvoterIDs(postID string, accountID string, exist bool)
	GetRoomList() []entity.Room
	GetRoom(id string) entity.Room
}
