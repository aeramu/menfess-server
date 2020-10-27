package service

type Repository interface{
	Save(Room) error
	FindByID(id string) (*Room, error)
	FindAll() (*[]Room, error)
	FindByStatus(status bool) (*[]Room, error)
}
