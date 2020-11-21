package service

type User struct{
	ID string
	Email string
	Password string
	Name string
	Avatar string
	Bio string
	PushToken map[string]bool
}
