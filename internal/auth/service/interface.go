package service

type Repository interface {
	Save(user User) error
	FindByEmail(email string) (*User, error)
}

type UserClient interface {
	Create() (*User, error)
}

type NotificationClient interface {
	AddPushToken(id string, pushToken string) error
	RemovePushToken(id string, pushToken string) error
}
