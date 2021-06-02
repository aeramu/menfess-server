package user

import (
	auth "github.com/aeramu/menfess-server/internal/auth/service"
	user "github.com/aeramu/menfess-server/internal/user/service"
)

func NewClient(user user.Service) auth.UserClient {
	return &userClient{user}
}

type userClient struct {
	user.Service
}

func (c *userClient) Create() (*auth.User, error) {
	u, err := c.Service.Create(user.CreateReq{
		Type: "user",
	})
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID: u.ID,
	}, nil
}
