package http

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"
)

type TokenClient struct {
	http.Client
	apiToken     *string
	isAuthorised bool
}

type transport struct {
	RoundTripper http.RoundTripper
	apiToken     *string
}

// NewTokenClient returns new http client with custom auth inherit
func NewTokenClient(timeout time.Duration, enableTLS bool) *TokenClient {
	token := ""

	var t transport
	// fix for using generated cert in dev
	if enableTLS {
		t = transport{
			RoundTripper: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig: &tls.Config{
					// UNSAFE! DON'T USE IN PRODUCTION!
					InsecureSkipVerify: true,
				},
			},
			apiToken: &token,
		}
	} else {
		t = transport{
			RoundTripper: http.DefaultTransport,
			apiToken:     &token,
		}
	}

	return &TokenClient{
		Client: http.Client{
			Timeout:   timeout,
			Transport: &t,
		},
		apiToken:     &token,
		isAuthorised: false,
	}
}

// SetToken sets auth token to client
func (c *TokenClient) SetToken(token string) {
	*c.apiToken = token
	c.isAuthorised = true
}

// DropAuth drops auth token to client
func (c *TokenClient) DropAuth() {
	*c.apiToken = ""
	c.isAuthorised = false
}

// DoWithAuth makes request with auth header
// returns error if client is not authenticated
func (c *TokenClient) DoWithAuth(req *http.Request) (*http.Response, error) {
	if len(*c.apiToken) == 0 {
		return nil, errors.New("client not authorized")
	}

	req.Header.Add("Authorization", *c.apiToken)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		c.DropAuth()
	}

	return resp, err
}

// RoundTrip Implements RoundTripper interface
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {

	return t.RoundTripper.RoundTrip(req)
}
