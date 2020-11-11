package service

type (
	RegisterReq struct{
		Email string
		Password string
		PushToken string
	}
	LoginReq struct {
		Email string
		Password string
		PushToken string
	}
	LogoutReq struct{
		ID string
		PushToken string
	}
	AuthReq struct{
		Token string
	}
)
