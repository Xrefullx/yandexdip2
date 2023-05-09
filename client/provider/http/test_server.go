package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// getTestHTTPServer returns test http server, server MUST be closed
// in response server returns headers, http status, body and wait sleep duration
// checks  require headers, body, url if not empty
func getTestHTTPServer(t *testing.T, cfg serverTestConfig) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			if len(cfg.reqMethod) > 0 {
				assert.EqualValues(t, cfg.reqMethod, req.Method)
			}

			if len(cfg.reqURL) > 0 {
				assert.EqualValues(t, cfg.reqURL, req.URL.Path)
			}

			//  check headers
			if len(cfg.reqHeaders) > 0 {
				for key, val := range cfg.reqHeaders {
					h := req.Header.Get(key)
					assert.EqualValues(t, val, h)
				}
			}

			//  check body
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			defer func() {
				if err := req.Body.Close(); err != nil {
					require.NoError(t, err)
				}
			}()

			assert.EqualValues(t, cfg.reqBody, string(body))

			//  return params
			if len(cfg.returnHeaders) > 0 {
				for key, val := range cfg.returnHeaders {
					rw.Header().Set(key, val)
				}
			}

			rw.WriteHeader(cfg.returnStatus)

			_, err = rw.Write([]byte(cfg.returnBody))
			require.NoError(t, err)

			if cfg.sleep != 0 {
				time.Sleep(cfg.sleep)
			}
		}))
}
