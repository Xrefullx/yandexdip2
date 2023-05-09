package http

import (
	"fmt"
	"github.com/Xrefullx/YanDip/client/provider"
	"io"
	"log"
	"net/http"
	"strings"
)

var _ provider.SecretProvider = (*HTTPProvider)(nil)

type HTTPProvider struct {
	client *TokenClient
	cfg    HTTPConfig
}

// NewHTTPProvider returns new http provider
// if base url starts with https, init client with tsl config
func NewHTTPProvider(cfg HTTPConfig) *HTTPProvider {
	withTLS := strings.HasPrefix(cfg.BaseURL, "https://")

	return &HTTPProvider{
		client: NewTokenClient(cfg.Timeout, withTLS),
		cfg:    cfg,
	}
}

// PingAuth checks connection with server and authentication
func (p *HTTPProvider) PingAuth() error {
	request, err := http.NewRequest(http.MethodGet, p.cfg.BaseURL+p.cfg.PingURL, nil)
	if err != nil {
		return fmt.Errorf("ping request error: %w", err)
	}

	//  do request
	response, err := p.client.DoWithAuth(request)
	if err != nil {
		return fmt.Errorf("ping request error: %w", err)
	}

	//  read body
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ping request error: wrong response: %v - %v", response.StatusCode, respBody)
	}
	return nil
}
