package service

type Payload struct {
	ID        string
	PushToken string
}

type User struct {
	ID       string
	Email    string
	Password string
}
