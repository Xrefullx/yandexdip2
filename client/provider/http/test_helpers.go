package http

import (
	"encoding/json"
	"github.com/Xrefullx/YanDip/client/provider/http/model"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	token = fake.CharactersN(16)

	authData = model.LoginRequest{
		Login:      fake.CharactersN(8),
		Password:   fake.CharactersN(8),
		MasterHash: fake.CharactersN(16),
		DeviceID:   uuid.New(),
	}

	srvBaseCfg = serverTestConfig{
		returnHeaders: map[string]string{"content-type": "application/json"},
		returnStatus:  http.StatusOK,
		sleep:         0,

		reqHeaders: map[string]string{"content-type": "application/json"},
	}

	provBaseCfg = HTTPConfig{
		AuthURL:     "/api/user/login",
		RegisterURL: "/api/user/register",
		SecretURL:   "/api/secret",
		SyncListURL: "/api/sync",
		PingURL:     "/api/ping",
		Timeout:     time.Millisecond * 500,
	}
)

func mustMarshal(val interface{}) string {
	res, err := json.Marshal(val)
	if err != nil {
		log.Fatal(err)
	}
	return string(res)
}

func getMockSyncList(count int) map[uuid.UUID]int {
	res := make(map[uuid.UUID]int)

	for i := 0; i < count; i++ {
		res[uuid.New()] = rand.Intn(10)
	}

	return res
}
