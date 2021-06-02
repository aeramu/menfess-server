package notification

import (
	"context"
	"encoding/json"
	"log"

	notif "github.com/aeramu/menfess-server/internal/notification/service"
	post "github.com/aeramu/menfess-server/internal/post/service"
	user "github.com/aeramu/menfess-server/internal/user/service"
)

func NewClient(notif notif.Service, user user.Service) post.NotificationClient {
	return &client{
		notif,
		user,
	}
}

type client struct {
	notif notif.Service
	user  user.Service
}

func (c *client) Send(event string, userID string, req interface{}) error {
	var title, body string
	data := map[string]string{}
	switch event {
	case "like":
		req := req.(post.LikeReq)
		u, err := c.user.Get(user.GetReq{ID: req.UserID})
		if err != nil {
			log.Println("Account Service Error:", err)
			return err
		}
		title = u.Name + " likes your post"
		data["id"] = req.PostID
	case "reply":
		p := req.(post.Post)
		u, err := c.user.Get(user.GetReq{ID: p.UserID})
		if err != nil {
			log.Println("Account Service Error:", err)
			return err
		}
		title = u.Name + " reply your post"
		body = p.Body
		data["id"] = p.ParentID
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = c.notif.SendNotification(context.Background(), notif.SendNotificationReq{
		Title:  title,
		Body:   body,
		UserID: userID,
		Data:   string(b),
	})
	if err != nil {
		return err
	}

	return nil
}
