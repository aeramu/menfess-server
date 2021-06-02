package service

//Repository interface
type Repository interface {
	Save(post Post) error
	FindByID(id string) (*Post, error)
	FindByParentID(id string, first int, after string, sort bool) ([]Post, error)
	FindByUserID(id string, first int, after string, sort bool) ([]Post, error)
	FindByAuthorID(id string, first int, after string, sort bool) ([]Post, error)
}

type NotificationClient interface {
	Send(event string, userID string, data interface{}) error
}
