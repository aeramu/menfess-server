package client

import (
	auth "github.com/aeramu/menfess-server/internal/auth/service"
	user "github.com/aeramu/menfess-server/internal/user/service"
)

func NewNotificationClient(user user.Service) auth.NotificationClient {
	return &userClient{user}
}

func (c *userClient) AddPushToken(id string, pushToken string) error {
	err := c.Service.AddPushToken(user.PushTokenReq{
		ID:        id,
		PushToken: pushToken,
	})
	return err
}

func (c *userClient) RemovePushToken(id string, pushToken string) error {
	err := c.Service.RemovePushToken(user.PushTokenReq{
		ID:        id,
		PushToken: pushToken,
	})
	return err
}
