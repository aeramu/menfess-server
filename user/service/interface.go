package service

type Repository interface {
	Save(user User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByType(t string) (*[]User, error)
}
