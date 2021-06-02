package service

type Repository interface {
	Save(user User) error
	FindByID(id string) (*User, error)
	FindByType(t string) ([]User, error)
}
