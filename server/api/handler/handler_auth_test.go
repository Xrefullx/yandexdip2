package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	apimodel "github.com/Xrefullx/YanDip/server/api/model"
	"github.com/Xrefullx/YanDip/server/model"
	mk "github.com/Xrefullx/YanDip/server/services/auth/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"log"
	"net/http"
	"testing"
)

var (
	mockUser = model.User{
		ID:           uuid.New(),
		Login:        fake.CharactersN(8),
		PasswordHash: fake.CharactersN(32),
	}

	mockLogin = apimodel.LoginRequest{
		Login:      mockUser.Login,
		Password:   fake.CharactersN(10),
		DeviceID:   uuid.New(),
		MasterHash: fake.CharactersN(32),
	}

	reqAuth = mustMarshalLogin(mockLogin)
)

func mustMarshalLogin(mockLogin apimodel.LoginRequest) string {
	res, err := json.Marshal(mockLogin)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(res)
}

// TestHandler_Login tests user register handler
func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []TestRoute{
		{
			name:    "return 200 if authenticated",
			method:  http.MethodPost,
			url:     "/api/user/login",
			svcAuth: authOk(ctrl),
			headers: map[string]string{"Content-Type": "application/json"},
			body:    reqAuth,

			expectedHeaders: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer tokentoken",
			},
			expectedCode: 200,
		},
		{
			name:         "return 500 if internal error",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authErrServer(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         reqAuth,
			expectedCode: 500,
		},
		{
			name:         "return 400 if wrong json format, missed quotes",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         fmt.Sprintf("{\"login: \"%s\",\"password\": \"%s\", \"device_id\": \"%s\"}", mockLogin.Login, mockLogin.Password, mockLogin.DeviceID.String()),
			expectedCode: 400,
		},
		{
			name:         "return 401 if wrong login/password",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authErrWrongPass(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         reqAuth,
			expectedCode: 401,
		},
		{
			name:         "return 415 if wrong content type",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			body:         reqAuth,
			expectedCode: 415,
		},
		{
			name:         "return 400 if empty login",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         fmt.Sprintf("{\"login\": \"%s\",\"password\": \"%s\", \"device_id\": \"%s\"}", "", mockLogin.Password, mockLogin.DeviceID.String()),
			expectedCode: 400,
		},
		{
			name:         "return 400 if empty password",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         fmt.Sprintf("{\"login\": \"%s\",\"password\": \"%s\", \"device_id\": \"%s\"}", mockLogin.Login, "", mockLogin.DeviceID.String()),
			expectedCode: 400,
		},
		{
			name:         "return 400 if empty device id",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         fmt.Sprintf("{\"login\": \"%s\",\"password\": \"%s\", \"device_id\": \"%s\"}", mockLogin.Login, mockLogin.Password, ""),
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}

}

// TestHandler_Register tests user login handler
func TestHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []TestRoute{
		{
			name:    "return 200 if registered",
			method:  http.MethodPost,
			url:     "/api/user/register",
			svcAuth: registerOk(ctrl),
			headers: map[string]string{"Content-Type": "application/json"},
			body:    reqAuth,

			expectedHeaders: map[string]string{"Content-Type": "application/json", "Authorization": "Bearer tokentoken"},
			expectedCode:    200,
		},
		{
			name:         "return 500 if internal error",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      registerErrServer(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         reqAuth,
			expectedCode: 500,
		},
		{
			name:         "return 409 if user exist",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      registerErrExist(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         reqAuth,
			expectedCode: 409,
		},
		{
			name:         "return 400 if wrong json format, missed quotes",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"iamuser\",\"password\": \"123456\"}",
			expectedCode: 400,
		},
		{
			name:         "return 415 if wrong content type",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			body:         reqAuth,
			expectedCode: 415,
		},
		{
			name:         "return 400 if empty login",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"\",\"password\": \"123456\"}",
			expectedCode: 400,
		},
		{
			name:         "return 400 if empty password",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"iamuser\",\"password\": \"\"}",
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}

}

/*  Auth mocks  */
func authEmpty(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	return authMock
}

/* Mocks for authenticate handler */
func authOk(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockUser, nil)
	authMock.EXPECT().EncodeTokenUserID(mockUser.ID, mockLogin.DeviceID, gomock.Any()).Return("tokentoken", nil)
	return authMock
}
func authErrWrongPass(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, model.ErrorWrongAuthData)
	return authMock
}
func authErrServer(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, errors.New("server error"))
	return authMock
}

/* Mocks for register handler */
func registerOk(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockUser, nil)
	authMock.EXPECT().EncodeTokenUserID(mockUser.ID, mockLogin.DeviceID, gomock.Any()).Return("tokentoken", nil)
	return authMock
}
func registerErrServer(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, errors.New("server error"))
	return authMock
}
func registerErrExist(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, model.ErrorConflictSaveUser)
	return authMock
}
