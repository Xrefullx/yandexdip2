package secret

import (
	"context"
	"fmt"
	"github.com/Xrefullx/YanDip/server/model"
	"github.com/Xrefullx/YanDip/server/storage"
	"github.com/google/uuid"
	"log"
)

var _ SecretManager = (*Secret)(nil)

type Secret struct {
	storage storage.Storage
}

func NewSecret(s storage.Storage) (*Secret, error) {
	return &Secret{
		storage: s,
	}, nil
}

func (s *Secret) AddUpdate(ctx context.Context, secret model.Secret) (uuid.UUID, int, error) {
	//  add if id is nil
	log.Printf("add secret %+v", secret)
	if secret.ID == uuid.Nil {
		if err := secret.ValidateAdd(); err != nil {
			return uuid.Nil, 0, fmt.Errorf("secret not valid to add: %w", err)
		}

		id, err := s.storage.Secret().Add(ctx, secret)
		if err != nil {
			return uuid.Nil, 0, err
		}

		return id, secret.Ver, nil
	}

	//  update
	if err := secret.ValidateUpdate(); err != nil {
		return uuid.Nil, 0, fmt.Errorf("secret not valid to update: %w", err)
	}

	dbSecret, err := s.storage.Secret().Get(ctx, secret.ID, secret.UserID)
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

	if err := s.storage.Secret().Update(ctx, dbSecret); err != nil {
		return uuid.Nil, 0, err
	}

	return dbSecret.ID, dbSecret.Ver, nil
}

func (s *Secret) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {

	if id == uuid.Nil {
		return fmt.Errorf("%w: id is nil", model.ErrorParamNotValid)
	}

	dbSecret, err := s.storage.Secret().Get(ctx, id, userID)
	if err != nil {
		return err
	}

	dbSecret.IsDeleted = true
	dbSecret.Ver = dbSecret.Ver + 1

	if err := s.storage.Secret().Update(ctx, dbSecret); err != nil {
		return err
	}

	return nil
}

func (s *Secret) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Secret, error) {
	if id == uuid.Nil {
		return model.Secret{}, fmt.Errorf("%w: id is nil", model.ErrorParamNotValid)
	}

	return s.storage.Secret().Get(ctx, id, userID)
}

func (s *Secret) GetUserSyncList(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("%w: userr id is nil", model.ErrorParamNotValid)
	}
	return s.storage.Secret().GetUserVersionList(ctx, userID)
}
