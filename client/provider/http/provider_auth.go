package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Xrefullx/YanDip/client/provider/http/model"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

func (p *HTTPProvider) Authorise(login string, pass string, masterHash string, deviceID uuid.UUID) error {
	return p.sendAuthorise(login, pass, masterHash, deviceID, p.cfg.BaseURL+p.cfg.AuthURL)
}

func (p *HTTPProvider) Register(login string, pass string, masterHash string, deviceID uuid.UUID) error {
	return p.sendAuthorise(login, pass, masterHash, deviceID, p.cfg.BaseURL+p.cfg.RegisterURL)
}

// Authorise make authorise request, get token and set it to client
func (p *HTTPProvider) sendAuthorise(login string, pass string, masterHash string, deviceID uuid.UUID, url string) error {
	loginData, err := json.Marshal(model.LoginRequest{
		Login:      login,
		Password:   pass,
		DeviceID:   deviceID,
		MasterHash: masterHash,
	})

	if err != nil {
		return err
	}

	//  prepare request
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(loginData))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	request.Header.Set("content-type", "application/json")

	//  do request
	response, err := p.client.Do(request)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
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

	// if not 200 error
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request error: response: %v - %s ", response.StatusCode, respBody)
	}

	//  get token
	token := response.Header.Get("Authorization")
	if len(token) == 0 {
		return fmt.Errorf("token is empty")
	}

	// set token to client
	p.client.SetToken(token)

	return nil
}
