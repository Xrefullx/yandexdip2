package middleware

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Middleware func(http.Handler) http.Handler

// HandlerTest stores test data and expected response params.
type TestMiddleware struct {
	name    string            //  test name
	method  string            //  http method
	body    string            //  request body content
	headers map[string]string //  request headers

	middlewareFunc []Middleware //  testing middleware func
	nextHandler    http.Handler // func for process request after testing middleware

	expectedBody    string            //  expected response body
	expectedHeaders map[string]string //  expected response headers
	expectedCode    int               // expected response http status
}

func (tt *TestMiddleware) CheckTest(t *testing.T) {

	h := tt.nextHandler
	for _, middleware := range tt.middlewareFunc {
		h = middleware(h)
	}

	request := httptest.NewRequest(tt.method, "/", bytes.NewBuffer([]byte(tt.body)))

	//  set headers
	if len(tt.headers) > 0 {
		for header, value := range tt.headers {
			request.Header.Set(header, value)
		}
	}

	//  make request
	w := httptest.NewRecorder()
	h.ServeHTTP(w, request)

	res := w.Result()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
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
