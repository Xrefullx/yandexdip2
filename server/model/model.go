package model

import (
	"fmt"
	"github.com/google/uuid"
)

type (
	User struct {
		ID           uuid.UUID
		Login        string
		PasswordHash string
		MasterHash   string
	}

	Secret struct {
		ID        uuid.UUID
		Ver       int
		UserID    uuid.UUID
		Data      string
		IsDeleted bool
	}
)

func (s *Secret) ValidateAdd() error {
	if s.ID != uuid.Nil {
		return fmt.Errorf("%w: id not nil", ErrorParamNotValid)
	}
	if s.Ver != 1 {
		return fmt.Errorf("%w: version not 0", ErrorParamNotValid)
	}
	if s.UserID == uuid.Nil {
		return fmt.Errorf("%w: user id is nil", ErrorParamNotValid)
	}
	if len(s.Data) == 0 {
		return fmt.Errorf("%w: data is empty", ErrorParamNotValid)
	}
	if s.IsDeleted {
		return fmt.Errorf("%w: is deleted", ErrorParamNotValid)
	}

	return nil
}

func (s *Secret) ValidateUpdate() error {
	if s.ID == uuid.Nil {
		return fmt.Errorf("%w: id is nil", ErrorParamNotValid)
	}
	if s.UserID == uuid.Nil {
		return fmt.Errorf("%w: user id is nil", ErrorParamNotValid)
	}
	if s.Ver == 0 {
		return fmt.Errorf("%w: version is 0", ErrorParamNotValid)
	}
	if !s.IsDeleted && len(s.Data) == 0 {
		return fmt.Errorf("%w: data is empty", ErrorParamNotValid)
	}

	return nil
}

func (u User) ValidateLogin(login string) error {
	if len(login) < 3 && len(login) > 60 {
		return fmt.Errorf("%w: login not valid", ErrorParamNotValid)
	}
	return nil
}
