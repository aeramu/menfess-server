package service

type (
	CreateReq struct {
		Email     string
		Password  string
		PushToken string
	}
	GetReq struct{
		ID string
	}
	GetByEmailReq struct{
		Email string
	}
	GetMenfessReq struct{}
	UpdateProfileReq struct {
		ID     string
		Name   string
		Avatar string
		Bio    string
	}
	PushTokenReq struct{
		ID string
		PushToken string
	}
)
