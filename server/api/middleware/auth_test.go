package middleware

import (
	"encoding/json"
	"github.com/Xrefullx/YanDip/server/api/model"
	"github.com/Xrefullx/YanDip/server/services/auth"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
)

// TestHandler_Login tests user register handler
func TestMiddlewareAuth(t *testing.T) {
	userdata := model.UserContextData{
		UserID:   uuid.New(),
		DeviceID: uuid.New(),
	}

	jsDataByte, err := json.Marshal(userdata)
	require.NoError(t, err)

	jsData := string(jsDataByte)

	//  jwtauth init
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	//  generate jwt token

	tokenString, err := auth.Auth{}.EncodeTokenUserID(userdata.UserID, userdata.DeviceID, tokenAuth)
	require.NoError(t, err)

	tokenNoUserIDString, err := auth.Auth{}.EncodeTokenUserID(uuid.Nil, userdata.DeviceID, tokenAuth)
	require.NoError(t, err)

	tokenNoDeviceIDString, err := auth.Auth{}.EncodeTokenUserID(userdata.UserID, uuid.Nil, tokenAuth)
	require.NoError(t, err)

	//  middleware reads Authorization header or jwt cookie to context
	jwtMiddleware := jwtauth.Verifier(tokenAuth)
	//  handler writes user_id to body
	toBodyHandler := writeUserIDToBody{t: t}

	tests := []TestMiddleware{
		{
			name:           "authenticated ok POST",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer " + tokenString},

			expectedBody: jsData,
			expectedCode: 200,
		},
		{
			name:           "authenticated ok GET",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodGet,
			headers:        map[string]string{"Authorization": "Bearer " + tokenString},

			expectedBody: jsData,
			expectedCode: 200,
		},
		{
			name:           "no jwt middleware - 401",
			middlewareFunc: []Middleware{MiddlewareAuth},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer " + tokenString},

			expectedCode: 401,
		},
		{
			name:           "no Authorization header - 401",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,

			expectedCode: 401,
		},
		{
			name:           "wrong token - 401",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer tokentoken"},

			expectedCode: 401,
		},

		{
			name:           "empty userID token - 401",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer " + tokenNoUserIDString},

			expectedCode: 401,
		},
		{
			name:           "no deviceID token - 401",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer " + tokenNoDeviceIDString},

			expectedCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}
}

// writeUserIDToBody handler with testing.T writes user_id from context to body
type writeUserIDToBody struct {
	t *testing.T
}

func (wr writeUserIDToBody) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value(model.ContextKeyUserID).(model.UserContextData)
	//require.NotEmpty(wr.t, ctxData)

	//userID, err := uuid.Parse(ctxID)
	//require.NoError(wr.t, err)

	data, err := json.Marshal(ctxData)
	require.NoError(wr.t, err)

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		log.Fatal(err.Error())
	}
}
