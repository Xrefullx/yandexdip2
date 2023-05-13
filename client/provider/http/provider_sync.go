package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/Xrefullx/YanDip/client/provider/http/model"
)

// GetSyncList downloads files meta info from server
func (p *HTTPProvider) GetSyncList() (map[uuid.UUID]int, error) {

	request, err := http.NewRequest(http.MethodGet, p.cfg.BaseURL+p.cfg.SyncListURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	request.Header.Set("content-type", "application/json")

	//  do request
	response, err := p.client.DoWithAuth(request)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	//  read body
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	var respObj model.SyncResponse
	if err := json.Unmarshal(respBody, &respObj); err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return respObj.List, nil
}
