package http

import (
	"github.com/Xrefullx/YanDip/client/provider/http/model"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestProviderSync_GetList(t *testing.T) {
	okList := getMockSyncList(40)
	token := fake.CharactersN(16)

	respData := model.SyncResponse{
		List: okList,
	}

	syncSrvConfig := srvBaseCfg.New(
		withReqMethod(http.MethodGet),
		withReqURL(provBaseCfg.SyncListURL),
	)

	tests := []struct {
		name      string
		serverCfg serverTestConfig

		reqErr  assert.ErrorAssertionFunc
		reqList map[uuid.UUID]int
		wait    time.Duration
	}{
		{
			name: "get exist list",
			serverCfg: syncSrvConfig.New(
				withReturnBody(mustMarshal(respData))),
			reqErr:  assert.NoError,
			reqList: okList,
		},
		{
			name: "too long request",
			serverCfg: syncSrvConfig.New(
				withReturnBody(mustMarshal(respData)),
				withSleep(provBaseCfg.Timeout+time.Millisecond*200)),
			reqErr: assert.Error,
		},
		{
			name: "err 500",
			serverCfg: syncSrvConfig.New(
				withReturnStatus(http.StatusInternalServerError)),
			reqErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := getTestHTTPServer(t, tt.serverCfg)
			defer server.Close()

			provCfg := provBaseCfg
			provCfg.BaseURL = server.URL

			provider := NewHTTPProvider(provCfg)
			provider.client.SetToken(token)
			list, err := provider.GetSyncList()

			tt.reqErr(t, err)

			//  check client authorised
			require.EqualValues(t, tt.reqList, list)
		})
	}
}
