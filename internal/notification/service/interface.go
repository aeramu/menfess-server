package service

import (
	"context"
)

//go:generate mockery --all

type Repository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user User) error
}

type PushServiceClient interface {
	Send(ctx context.Context, pushToken map[string]bool, title string, body string, data string) error
}
