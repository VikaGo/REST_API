package service

import (
	"context"

	"github.com/VikaGo/REST_API/store"
	"github.com/pkg/errors"
)

// Manager is just a collection of all services we have in the project
type Manager struct {
	User UserService
}

// NewManager creates new service manager
func NewManager(ctx context.Context, store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No store provided")
	}
	return &Manager{
		User: NewUserWebService(ctx, store),
	}, nil
}
