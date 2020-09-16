package usecase

import "github.com/aeramu/menfess-server/entity"

//Interactor interface
type Interactor interface {
	post
	room
}

type post interface {
	Post(id string) entity.Post
	PostFeed(first int, after string) []entity.Post
	PostChild(parentID string, first int, after string) []entity.Post
	PostRooms(roomIDs []string, first int, after string) []entity.Post
	PostPost(name string, avatar string, body string, parentID string, repostID string, roomID string) entity.Post
	UpvotePost(accountID string, postID string) entity.Post
	DownvotePost(accountID string, postID string) entity.Post
}

type room interface {
	RoomList() []entity.Room
	Room(id string) entity.Room
}

//InteractorConstructor constructor
type InteractorConstructor struct {
	Repository Repository
}

//New Construct Interactor
func (i InteractorConstructor) New() Interactor {
	return &interactor{
		repo: i.Repository,
	}
}
