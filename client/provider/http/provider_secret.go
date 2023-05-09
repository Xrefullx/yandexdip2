package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xrefullx/YanDip/client/model"
	prmodel "github.com/Xrefullx/YanDip/client/provider/http/model"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

// UploadSecret uploads secret to server, returns server id and version
// if id is nil, creates new
func (p *HTTPProvider) UploadSecret(data string, id uuid.UUID, ver int) (uuid.UUID, int, error) {
	reqData := prmodel.SecretRequest{
		Data: data,
		ID:   id,
		Ver:  ver,
	}
	if err := reqData.ValidateUpload(); err != nil {
		return uuid.Nil, 0, err
	}

	resp := prmodel.SecretRequest{}
	if err := p.processSecretRequest(reqData, http.MethodPut, &resp); err != nil {
		return uuid.Nil, 0, fmt.Errorf("error secret upload: %w", err)
	}

	if !resp.IsValidResponseUpload() {
		return uuid.Nil, 0, errors.New("upload response error: response not valid")
	}

	return resp.ID, resp.Ver, nil
}

// DownloadSecret downloads secret from server
func (p *HTTPProvider) DownloadSecret(id uuid.UUID) (uuid.UUID, int, string, error) {
	reqData := prmodel.SecretRequest{
		ID: id,
	}
	if !reqData.IsValidDownload() {
		return uuid.Nil, 0, "", fmt.Errorf("%w : not valid download param", model.ErrorParamNotValid)
	}

	resp := prmodel.SecretRequest{}
	if err := p.processSecretRequest(reqData, http.MethodGet, &resp); err != nil {
		return uuid.Nil, 0, "", fmt.Errorf("error secret download: %w", err)
	}

	if !resp.IsValidResponseDownload() {
		return uuid.Nil, 0, "", errors.New("download response error: response not valid")
	}
	return resp.ID, resp.Ver, resp.Data, nil
}

// DeleteSecret deletes secret from server
func (p *HTTPProvider) DeleteSecret(id uuid.UUID) error {
	reqData := prmodel.SecretRequest{
		ID: id,
	}
	if !reqData.IsValidDelete() {
		return fmt.Errorf("%w : not valid delete param", model.ErrorParamNotValid)
	}

	resp := prmodel.SecretRequest{}
	if err := p.processSecretRequest(reqData, http.MethodDelete, &resp); err != nil {
		return fmt.Errorf("error secret delete: %w", err)
	}
	return nil
}

func (p *HTTPProvider) processSecretRequest(req prmodel.SecretRequest, method string, resp *prmodel.SecretRequest) error {
	secretData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("secret %v request error: %w", method, err)
	}

	//  prepare request
	request, err := http.NewRequest(method, p.cfg.BaseURL+p.cfg.SecretURL, bytes.NewBuffer(secretData))
	if err != nil {
		return fmt.Errorf("secret %v request error: %w", method, err)
	}

	request.Header.Set("content-type", "application/json")

	//  do request
	response, err := p.client.DoWithAuth(request)
	if err != nil {
		return fmt.Errorf("secret %v request error: %w", method, err)
	}

	//  read body
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("secret %v response error: %w", method, err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("secret %s response error: %v - %s", method, response.StatusCode, string(respBody))
	}

	//  if body not empty unmarshal
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, resp); err != nil {
			return fmt.Errorf("secret %v response error: %w body: %v ", method, err, respBody)
		}
	}

	return nil
}
