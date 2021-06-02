package service

type User struct {
	ID        string
	PushToken map[string]bool
}

func (u *User) AddPushToken(pushToken string) {
	if u.PushToken == nil {
		u.PushToken = map[string]bool{}
	}
	u.PushToken[pushToken] = true
}

func (u *User) RemovePushToken(pushToken string) {
	delete(u.PushToken, pushToken)
}
