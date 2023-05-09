package handler

import (
	"bytes"
	"github.com/Xrefullx/YanDip/server/services/auth"
	"github.com/Xrefullx/YanDip/server/services/secret"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http/httptest"
	"testing"
)

// HandlerTest stores test data and expected response params.
type TestRoute struct {
	svcAuth   auth.Authenticator //  authentication service
	svcSecret secret.SecretManager
	jwtAuth   *jwtauth.JWTAuth

	name    string            //  test name
	method  string            //  http method
	url     string            //  request url
	body    string            //  request body content
	headers map[string]string //  request headers

	expectedBody    string            //  expected response body
	expectedHeaders map[string]string //  expected response headers
	expectedCode    int               // expected response http status
}

// CheckTest runs handler, builds request and checks response values
func (tt *TestRoute) CheckTest(t *testing.T) {
	//  new handler with mock services
	h, err := NewHandler(tt.svcAuth, tt.svcSecret, tt.jwtAuth)
	require.NoError(t, err)

	//  new router with handler
	r := GetRouter(h)
	request := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte(tt.body)))

	//  set headers
	if len(tt.headers) > 0 {
		for header, value := range tt.headers {
			request.Header.Set(header, value)
		}
	}

	//  make request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	res := w.Result()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	strBody := string(resBody)

	//  check response headers
	if len(tt.expectedHeaders) > 0 {
		for expHeader, expValue := range tt.expectedHeaders {
			resHeader := res.Header.Get(expHeader)
			assert.Equal(t, expValue, resHeader, "ожидался заголовок %v, значение %v, получено значение %v",
				expHeader, expValue, resHeader)
		}
	}

	//  check response code
	if tt.expectedCode != 0 {
		assert.True(t, res.StatusCode == tt.expectedCode, "Ожидался код ответа %d, получен %d", tt.expectedCode, w.Code)
	}

	//  check response body
	if tt.expectedBody != "" {
		assert.Equal(t, strBody, tt.expectedBody, "Ожидался ответа %v, получен %v", tt.expectedBody, strBody)
	}
}
