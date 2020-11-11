package service

type UserClient interface {
	Create(email string, password string, pushToken string) (*User, error)
	GetByEmail(email string) (*User, error)
	AddPushToken(id string, pushToken string) error
	RemovePushToken(id string, pushToken string) error
}
