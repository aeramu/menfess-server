package client

import (
	post "github.com/aeramu/menfess-server/post/service"
	user "github.com/aeramu/menfess-server/user/service"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"log"
)

func NewNotificationClient(user user.Service) post.NotificationClient{
	return &notificationClient{
		expo.NewPushClient(nil),
		user,
	}
}

type notificationClient struct {
	client *expo.PushClient
	user user.Service
}

func (c *notificationClient) Send(event string, userID string, req interface{}) error {
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

	u, err := c.user.Get(user.GetReq{ID: userID})
	if err != nil {
		log.Println("Account Service Error:", err)
		return err
	}
	
	var tokens []expo.ExponentPushToken
	for t := range u.PushToken{
		token, err := expo.NewExponentPushToken(t)
		if err != nil{
			continue
		}
		tokens = append(tokens, token)
	}
	_, err = c.client.Publish(&expo.PushMessage{
		To:    tokens,
		Body:  body,
		Data:  data,
		Sound: "default",
		Title: title,
	})
	if err != nil{
		return err
	}
	return nil
}

