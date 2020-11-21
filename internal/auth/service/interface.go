package service

//go:generate mockery --all --keeptree --case underscore

type UserClient interface {
	Create(email string, password string, pushToken string) (*User, error)
	GetByEmail(email string) (*User, error)
}

type NotificationClient interface {
	AddPushToken(id string, pushToken string) error
	RemovePushToken(id string, pushToken string) error
}
