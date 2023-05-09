package secret

import (
	"context"
	"github.com/Xrefullx/YanDip/server/model"
	"github.com/google/uuid"
)

type SecretManager interface {
	AddUpdate(ctx context.Context, secret model.Secret) (uuid.UUID, int, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Secret, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetUserSyncList(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int, error)
}
