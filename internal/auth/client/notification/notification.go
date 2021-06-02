package notification

import (
	"context"

	auth "github.com/aeramu/menfess-server/internal/auth/service"
	notif "github.com/aeramu/menfess-server/internal/notification/service"
)

func NewClient(notif notif.Service) auth.NotificationClient {
	return &client{notif}
}

type client struct {
	notif notif.Service
}

func (c *client) AddPushToken(id string, pushToken string) error {
	err := c.notif.AddPushToken(context.Background(), notif.AddPushTokenReq{
		ID:        id,
		PushToken: pushToken,
	})
	return err
}

func (c *client) RemovePushToken(id string, pushToken string) error {
	err := c.notif.RemovePushToken(context.Background(), notif.RemovePushTokenReq{
		ID:        id,
		PushToken: pushToken,
	})
	return err
}
