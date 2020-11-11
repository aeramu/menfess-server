package client

import (
	auth "github.com/aeramu/menfess-server/auth/service"
	user "github.com/aeramu/menfess-server/user/service"
)

func NewUserClient(user user.Service) auth.UserClient{
	return &userClient{user}
}

type userClient struct{
	user.Service
}

func (c *userClient) Create(email string, password string, pushToken string) (*auth.User, error) {
	u, err := c.Service.Create(user.CreateReq{
		Email:     email,
		Password:  password,
		PushToken: pushToken,
	})
	if err != nil{
		return nil, err
	}
	return &auth.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (c *userClient) GetByEmail(email string) (*auth.User, error) {
	u, err := c.Service.GetByEmail(user.GetByEmailReq{Email: email})
	if err != nil{
		return nil, err
	}
	if u == nil{
		return nil, nil
	}
	return &auth.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, nil
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

