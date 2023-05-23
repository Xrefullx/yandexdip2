package secret

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Xrefullx/YanDip/server/model"
	"github.com/Xrefullx/YanDip/server/storage"
)

type Secret struct {
	storage storage.SecretRepository
}

func NewSecret(repo storage.SecretRepository) (*Secret, error) {
	return &Secret{
		storage: repo,
	}, nil
}

func (s *Secret) Add(ctx context.Context, secret model.Secret) (uuid.UUID, int, error) {
	log.Printf("add secret %+v", secret)
	if err := secret.ValidateAdd(); err != nil {
		return uuid.Nil, 0, fmt.Errorf("secret not valid to add: %w", err)
	}

	id, err := s.storage.Add(ctx, secret)
	if err != nil {
		return uuid.Nil, 0, err
	}

	return id, secret.Ver, nil
}
func (s *Secret) Update(ctx context.Context, secret model.Secret) (uuid.UUID, int, error) {
	if err := secret.ValidateUpdate(); err != nil {
		return uuid.Nil, 0, fmt.Errorf("secret not valid to update: %w", err)
	}

	dbSecret, err := s.storage.Get(ctx, secret.ID, secret.UserID)
	if err != nil {
		return uuid.Nil, 0, err
	}

	if dbSecret.IsDeleted {
		return uuid.Nil, 0, model.ErrorItemIsDeleted
	}

	//  incoming version can be > local, in collision fix
	//  version increments by server version
	if dbSecret.Ver > secret.Ver {
		return uuid.Nil, 0, model.ErrorVersionToLow
	}

	dbSecret.Data = secret.Data
	dbSecret.Ver = dbSecret.Ver + 1

	if err := s.storage.Update(ctx, dbSecret); err != nil {
		return uuid.Nil, 0, err
	}

	return dbSecret.ID, dbSecret.Ver, nil
}

func (s *Secret) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: id is nil", model.ErrorParamNotValid)
	}

	return s.storage.Delete(ctx, id, userID)
}

func (s *Secret) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Secret, error) {
	if id == uuid.Nil {
		return model.Secret{}, fmt.Errorf("%w: id is nil", model.ErrorParamNotValid)
	}

	return s.storage.Get(ctx, id, userID)
}

func (s *Secret) GetUserSyncList(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("%w: userr id is nil", model.ErrorParamNotValid)
	}
	return s.storage.GetUserVersionList(ctx, userID)
}
