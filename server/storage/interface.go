package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/Xrefullx/YanDip/server/model"
)

type UserRepository interface {
	//  Adds new user to storage
	Create(ctx context.Context, user model.User) (model.User, error)
	//  Returns user from storage
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

type SecretRepository interface {
	Add(ctx context.Context, secret model.Secret) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Secret, error)
	Update(ctx context.Context, secret model.Secret) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	GetUserVersionList(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int, error)
}
