package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Xrefullx/YanDip/server/model"
	"github.com/Xrefullx/YanDip/server/storage"
)

var _ Authenticator = (*Auth)(nil)

const (
	salt = "nJkksjjdxszx120_dssd!xc"
)

// Auth implements Authenticator interface methods for user authorisation.
type Auth struct {
	userRepo storage.UserRepository
}

// NewAuth init new Auth.
func NewAuth(repo storage.UserRepository) (*Auth, error) {
	return &Auth{
		userRepo: repo,
	}, nil
}

// CreateUser creates new user.
// If user exist, returns ErrorUserAlreadyExist.
func (a *Auth) CreateUser(ctx context.Context, login string, password string, masterHash string) (model.User, error) {
	loginHash, err := bcrypt.GenerateFromPassword([]byte(password+salt), 10)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", model.ErrAddingUser, err)
	}

	masterHashHash, err := bcrypt.GenerateFromPassword([]byte(masterHash+salt), 10)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", model.ErrAddingUser, err)
	}

	user := model.User{
		Login:        login,
		PasswordHash: string(loginHash),
		MasterHash:   string(masterHashHash),
	}

	user, err = a.userRepo.Create(ctx, user)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %v", model.ErrAddingUser, err)
	}

	return user, err
}

// Authenticate checks user login, password and return.
// If user not founded, or wrong password, returns ErrorWrongAuthData.
func (a *Auth) Authenticate(ctx context.Context, login string, password string, masterHash string) (model.User, error) {
	user, err := a.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrorItemNotFound) {
			return model.User{}, model.ErrorWrongAuthData
		}

		return model.User{}, fmt.Errorf("%w: %v", model.ErrAuthenticatingUser, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password+salt)); err != nil {
		return model.User{}, model.ErrorWrongAuthData
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(masterHash+salt)); err != nil {
		return model.User{}, model.ErrorWrongAuthData
	}

	return user, nil
}

// EncodeTokenUserID encodes token with user_id claim.
func (a Auth) EncodeTokenUserID(userID uuid.UUID, deviceID uuid.UUID, tokenAuth *jwtauth.JWTAuth) (string, error) {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id":   userID.String(),
		"device_id": deviceID.String(),
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", model.ErrGeneratingToken, err)
	}

	return tokenString, nil
}
