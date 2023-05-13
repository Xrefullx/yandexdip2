package storage

import (
	"github.com/google/uuid"

	"github.com/Xrefullx/YanDip/client/model"
)

type Storage interface {
	AddSecret(v model.Secret) (int64, error)
	GetSecret(id int64) (model.Secret, error)
	GetSecretByExtID(extID uuid.UUID) (model.Secret, error)
	GetMetaList() ([]model.SecretMeta, error)

	//UpdateSecretBySecretID(v model.Secret) error
	UpdateSecret(v model.Secret) error
	DeleteSecret(id int64) error
	Close()
}
