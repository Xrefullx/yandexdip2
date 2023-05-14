package model

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type (
	User struct {
		ID           uuid.UUID
		Login        string `validate:"required,min=3,max=60"`
		PasswordHash string `validate:"required"`
		MasterHash   string `validate:"required"`
	}

	Secret struct {
		ID        uuid.UUID `validate:"required"`
		Ver       int       `validate:"required,min=1"`
		UserID    uuid.UUID `validate:"required"`
		Data      string    `validate:"required_without=IsDeleted"`
		IsDeleted bool
	}
)

func (s *Secret) ValidateAdd() error {
	if s.ID != uuid.Nil {
		return fmt.Errorf("%w: id not nil", ErrorParamNotValid)
	}

	err := validate.Struct(s)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorParamNotValid, err)
	}

	return nil
}

func (s *Secret) ValidateUpdate() error {
	if s.ID == uuid.Nil {
		return fmt.Errorf("%w: id is nil", ErrorParamNotValid)
	}

	err := validate.Struct(s)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorParamNotValid, err)
	}

	return nil
}

func (u *User) ValidateLogin() error {
	err := validate.Var(u.Login, "required,min=3,max=60")
	if err != nil {
		return fmt.Errorf("%w: login not valid", ErrorParamNotValid)
	}

	return nil
}
