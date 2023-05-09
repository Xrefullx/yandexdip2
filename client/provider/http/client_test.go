package http

import (
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenClient_SetToken(t *testing.T) {
	client := NewTokenClient(time.Second*5, false)

	token := fake.CharactersN(32)
	client.SetToken(token)
	testTokenClient(t, client, token)
}

func testTokenClient(t *testing.T, client *TokenClient, reqToken string) {
	server := httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			rToken := req.Header.Get("Authorization")
			log.Println(rToken)
			if _, err := rw.Write([]byte(rToken)); err != nil {
				log.Fatal(err.Error())
			}
		}))
	defer server.Close()

	request, err := http.NewRequest(http.MethodGet, server.URL, nil)

	resp, err := client.DoWithAuth(request)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	assert.EqualValues(t, reqToken, string(body))
}
