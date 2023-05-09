package services

import (
	"github.com/google/uuid"
)

var (
	SyncActions = map[string]int{
		"UPLOAD_NEW":     1,
		"UPLOAD":         2,
		"DOWNLOAD_NEW":   3,
		"DOWNLOAD":       4,
		"DELETE_LOCALLY": 5,
		"SEND_DELETE":    6,
		"COLLISION":      7,
	}
)

type SyncTask struct {
	LocID     int64
	SecretId  uuid.UUID
	Ver       int
	ActionID  int
	TimeStamp int64
}
