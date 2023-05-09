package auth

import (
	"context"
	"github.com/Xrefullx/YanDip/server/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

// Authenticator is the interface that wraps methods user identification, authentication, authorisation.
type Authenticator interface {
	CreateUser(ctx context.Context, login string, password string, masterHash string) (model.User, error)
	Authenticate(ctx context.Context, login string, password string, masterHash string) (model.User, error)
	EncodeTokenUserID(userID uuid.UUID, deviceID uuid.UUID, tokenAuth *jwtauth.JWTAuth) (string, error)
}
