package model

import (
	"fmt"
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Login      string    `json:"login"`
	MasterHash string    `json:"master_hash"`
	Password   string    `json:"password"`
	DeviceID   uuid.UUID `json:"device_id"`
}

type SecretRequest struct {
	Data string    `json:"data,omitempty"`
	ID   uuid.UUID `json:"id,omitempty"`
	Ver  int       `json:"ver,omitempty"`
}

func (s *SecretRequest) IsValidResponseUpload() bool {
	return s.ID != uuid.Nil && s.Ver > 0
}

func (s *SecretRequest) IsValidResponseDownload() bool {
	return s.ID != uuid.Nil && s.Ver > 0 && len(s.Data) > 0
}

func (s *SecretRequest) ValidateUpload() error {
	if s.ID == uuid.Nil && s.Ver > 1 {
		return fmt.Errorf("%w: ver is %v and secret id is nil", model.ErrorParamNotValid, s.Ver)
	}
	if s.Ver == 0 {
		return fmt.Errorf("%w: ver is %v", model.ErrorParamNotValid, s.Ver)
	}
	if len(s.Data) == 0 {
		return fmt.Errorf("%w: data is empty", model.ErrorParamNotValid)
	}
	return nil
}
func (s *SecretRequest) IsValidDelete() bool {
	if s.ID == uuid.Nil {
		return false
	}
	if len(s.Data) != 0 || s.Ver > 0 {
		return false
	}
	return true
}

func (s *SecretRequest) IsValidDownload() bool {
	if s.ID == uuid.Nil {
		return false
	}
	if len(s.Data) != 0 || s.Ver > 0 {
		return false
	}
	return true
}

type SyncResponse struct {
	List map[uuid.UUID]int `json:"list"`
}
