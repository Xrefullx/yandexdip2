package services

import (
	"errors"
	"fmt"
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/google/uuid"
	"log"
)

// uploads secret to server and updates local SecretID and version
// If UPLOAD_NEW - task exist only local, get by id
// If UPLOAD - task was synch and have SecretID, get by secret id
// If response 200, write secret meta data from response and set status ACTUAL
func (s *SyncService) Upload(task SyncTask) error {
	var secret model.Secret
	var err error

	if task.ActionID == SyncActions["UPLOAD_NEW"] {
		secret, err = s.db.GetSecret(task.LocID)
	}

	if task.ActionID == SyncActions["UPLOAD"] {
		secret, err = s.db.GetSecretByExtID(task.SecretId)
	}

	if err != nil {
		return err
	}

	id, ver, err := s.provider.UploadSecret(secret.SecretData, task.SecretId, task.Ver)
	if err != nil {
		return err
	}

	if secret.SecretID != uuid.Nil && secret.SecretID != id {
		return errors.New("error upload sync: response secretID not equal local")
	}

	secret.SecretVer = ver
	secret.SecretID = id
	secret.StatusID = model.SecretStatuses["ACTUAL"]

	if err := s.db.UpdateSecret(secret); err != nil {
		return fmt.Errorf("error upload sync: error save secret meta info: %w", err)
	}

	return nil
}

// Download downloads secret from server
// If response 200, updates local data and meta.
func (s *SyncService) Download(task SyncTask) error {
	id, ver, data, err := s.provider.DownloadSecret(task.SecretId)
	if err != nil {
		return err
	}

	info := model.Info{}
	if err := info.FromEncodedData(data, s.cfg.MasterKey); err != nil {
		return fmt.Errorf("error read info drom encoded secret data: %w", err)
	}

	dbSecret, err := s.db.GetSecretByExtID(id)
	if err != nil {
		return fmt.Errorf("error save secret data to storage: %w", err)
	}

	dbSecret.Info = info
	dbSecret.SecretVer = ver
	dbSecret.StatusID = model.SecretStatuses["ACTUAL"]
	dbSecret.SecretData = data

	if err := s.db.UpdateSecret(dbSecret); err != nil {
		return fmt.Errorf("error save secret data to storage: %w", err)
	}

	return nil
}

// DownloadNew downloads secret from server
// If response 200, creates new local data.
func (s *SyncService) DownloadNew(task SyncTask) error {
	id, ver, data, err := s.provider.DownloadSecret(task.SecretId)
	if err != nil {
		return err
	}

	info := model.Info{}
	if err := info.FromEncodedData(data, s.cfg.MasterKey); err != nil {
		return fmt.Errorf("error read info drom encoded secret data: %w", err)
	}

	_, err = s.db.AddSecret(model.Secret{
		Info:       info,
		SecretID:   id,
		SecretVer:  ver,
		StatusID:   model.SecretStatuses["ACTUAL"],
		SecretData: data,
	})

	if err != nil {
		return fmt.Errorf("error save secret data to storage: %w", err)
	}

	return nil
}

// DeleteRemote deletes secret from server
// If response 200, mark local secret status as DELETED.
func (s *SyncService) DeleteRemote(task SyncTask) error {
	secret, err := s.db.GetSecretByExtID(task.SecretId)
	if err != nil {
		return err
	}

	if err := s.provider.DeleteSecret(task.SecretId); err != nil {
		return err
	}

	secret.StatusID = model.SecretStatuses["DELETED"]

	err = s.db.UpdateSecret(secret)
	if err != nil {
		return err
	}

	return nil
}

// DeleteLocally deletes secret from local database
func (s *SyncService) DeleteLocally(task SyncTask) error {
	return s.db.DeleteSecret(task.LocID)
}

// DeleteLocally deletes secret from local database
func (s *SyncService) AddCollision(task SyncTask) error {
	// use external version - delete local, then ext version downloads on sync
	// use local version - set loc ver as external
	return nil
}

func (s *SyncService) ProcessTask(task SyncTask) error {
	log.Printf("task started %+v", task)
	switch task.ActionID {

	case SyncActions["UPLOAD"], SyncActions["UPLOAD_NEW"]:
		if err := s.Upload(task); err != nil {
			return err
		}
	case SyncActions["DOWNLOAD"]:
		if err := s.Download(task); err != nil {
			return err
		}
	case SyncActions["DOWNLOAD_NEW"]:
		if err := s.DownloadNew(task); err != nil {
			return err
		}
	case SyncActions["SEND_DELETE"]:
		if err := s.DeleteRemote(task); err != nil {
			return err
		}
	case SyncActions["DELETE_LOCALLY"]:
		if err := s.DeleteLocally(task); err != nil {
			return err
		}
	case SyncActions["COLLISION"]:
		if err := s.AddCollision(task); err != nil {
			return err
		}
	}

	log.Printf("task processed %+v", task)

	return nil
}
