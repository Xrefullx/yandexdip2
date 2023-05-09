package provider

import (
	"github.com/google/uuid"
)

type SecretProvider interface {
	Authorise(login string, pass string, masterHash string, deviceID uuid.UUID) error
	Register(login string, pass string, masterHash string, deviceID uuid.UUID) error
	PingAuth() error

	UploadSecret(data string, id uuid.UUID, ver int) (uuid.UUID, int, error)
	DownloadSecret(id uuid.UUID) (uuid.UUID, int, string, error)
	DeleteSecret(id uuid.UUID) error
	GetSyncList() (map[uuid.UUID]int, error)
}
