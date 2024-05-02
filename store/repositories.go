package store

import (
	"context"

	"github.com/VikaGo/REST_API/model"
	"github.com/google/uuid"
)

// UserRepo is a store for users
//
//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	GetUser(context.Context, uuid.UUID) (*model.DBUser, error)
	CreateUser(context.Context, *model.DBUser) (*model.DBUser, error)
	UpdateUser(context.Context, *model.DBUser) (*model.DBUser, error)
	DeleteUser(context.Context, uuid.UUID) error
	GetPassword(ctx context.Context, id uuid.UUID) (string, error)
	GetUserByNickname(ctx context.Context, nickname string) (*model.DBUser, error)
}
