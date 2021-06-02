package service

type (
	AddPushTokenReq struct {
		ID        string
		PushToken string
	}
	RemovePushTokenReq struct {
		ID        string
		PushToken string
	}
	SendNotificationReq struct {
		Title  string
		Body   string
		UserID string
		Data   string
	}
)
