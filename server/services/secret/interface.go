package secret

import (
	"context"

	"github.com/google/uuid"

	"github.com/Xrefullx/YanDip/server/model"
)

type SecretManager interface {
	Add(ctx context.Context, secret model.Secret) (uuid.UUID, int, error)
	Update(ctx context.Context, secret model.Secret) (uuid.UUID, int, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Secret, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Secret, error)
	GetUserSyncList(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int, error)
}
