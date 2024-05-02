package service

import (
	"context"

	"github.com/VikaGo/REST_API/model"
	"github.com/google/uuid"
)

type UserService interface {
	GetUser(context.Context, uuid.UUID) (*model.User, error)
	CreateUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
	DeleteUser(context.Context, uuid.UUID) error
	GetPassword(context.Context, uuid.UUID) (string, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error
	GetUserByNickname(ctx context.Context, nickname string) (*model.User, error)
	GenerateToken(ctx context.Context, nickname string, password string) (string, error)
}
