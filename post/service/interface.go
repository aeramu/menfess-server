package service

//Repository interface
type Repository interface {
	Save(post Post) error
	FindByID(id string) (*Post, error)
	FindByParentID(id string, first int, after string, sort bool) (*[]Post, error)
	FindByAuthorID(id string, first int, after string, sort bool) (*[]Post, error)
	//GetPostByID(id string) Post
	//GetPostListByParentID(parentID string, first int, after string, ascSort bool) []Post
	//GetPostListByRoomIDs(roomIDs []string, first int, after string, ascSort bool) []Post
	//PutPost(name string, avatar string, body string, parentID string, repostID string, roomID string) Post
	//UpdateUpvoterIDs(postID string, accountID string, exist bool)
	//UpdateDownvoterIDs(postID string, accountID string, exist bool)
	//GetRoomList() []Room
	//GetRoom(id string) Room
}

type NotificationClient interface {
	Send(event string, userID string, data interface{}) error
}
