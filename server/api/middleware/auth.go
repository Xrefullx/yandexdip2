package middleware

import (
	"context"
	"errors"
	"github.com/Xrefullx/YanDip/server/api/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
)

// MiddlewareAuth gets token from request, checks it and sets user_id to context
func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if token == nil {
			http.Error(w, "token is nil", http.StatusUnauthorized)
			return
		}

		if err := jwt.Validate(token); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userData, err := readClaims(claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Set userID to context
		ctx := context.WithValue(r.Context(), model.ContextKeyUserID, userData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func readClaims(claims map[string]interface{}) (model.UserContextData, error) {
	objUserID, ok := claims["user_id"]
	if !ok {
		return model.UserContextData{}, errors.New("user id is empty")
	}
	userID, err := uuid.Parse(objUserID.(string))
	if err != nil || userID == uuid.Nil {
		return model.UserContextData{}, errors.New("wrong user id")
	}

	objDeviceID, ok := claims["device_id"]
	if !ok {
		return model.UserContextData{}, errors.New("device id is empty")
	}
	deviceID, err := uuid.Parse(objDeviceID.(string))
	if err != nil || deviceID == uuid.Nil {
		return model.UserContextData{}, errors.New("wrong device id")
	}

	return model.UserContextData{
		UserID:   userID,
		DeviceID: deviceID,
	}, nil
}
