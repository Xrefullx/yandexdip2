package services

import (
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/Xrefullx/YanDip/client/pkg"
	"github.com/Xrefullx/YanDip/client/storage"
	mk "github.com/Xrefullx/YanDip/client/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var cfg = pkg.Config{
	MasterKey: "testKey",
}

func GetTestSecretSvc(t *testing.T, storage storage.Storage) *SecretService {
	svcSecret := NewSecret(&cfg, storage)
	return &svcSecret
}

func TestCard_ToSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	binary, err := SecretService{}.ReadBinary("test_img.svg")
	require.NoError(t, err)

	tests := []struct {
		name    string
		obj     interface{}
		storage storage.Storage
		reqErr  assert.ErrorAssertionFunc
	}{
		{
			name:    "card",
			obj:     model.TestCard,
			storage: storageEmpty(ctrl),
		},
		{
			name:    "text",
			obj:     model.TestText,
			storage: storageEmpty(ctrl),
		},
		{
			name:    "auth",
			obj:     model.TestAuth,
			storage: storageEmpty(ctrl),
		},
		{
			name:    "binary",
			obj:     binary,
			storage: storageEmpty(ctrl),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			secretSvc := GetTestSecretSvc(t, tt.storage)

			secret, err := secretSvc.ToSecret(tt.obj)
			require.NoError(t, err)

			info := model.Info{}
			err = info.FromEncodedData(secret.SecretData, cfg.MasterKey)
			require.NoError(t, err)

			resObj, err := secretSvc.ReadFromSecret(secret)
			require.NoError(t, err)

			switch resObj.(type) {
			case model.Card:
				src := tt.obj.(model.Card)
				res := resObj.(model.Card)
				require.Equal(t, src.Info, info)
				require.Equal(t, tt.obj, res)
			case model.Auth:
				src := tt.obj.(model.Auth)
				res := resObj.(model.Auth)
				require.Equal(t, src.Info, info)
				require.Equal(t, tt.obj, res)
			case model.Text:
				src := tt.obj.(model.Text)
				res := resObj.(model.Text)
				require.Equal(t, src.Info, info)
				require.Equal(t, tt.obj, res)
			case model.Binary:
				src := tt.obj.(model.Binary)
				res := resObj.(model.Binary)
				require.Equal(t, src.Info, info)
				require.Equal(t, tt.obj, res)

			default:
				t.Error("wrong type")
			}

		})
	}
}

func storageEmpty(ctrl *gomock.Controller) *mk.MockStorage {
	storageMock := mk.NewMockStorage(ctrl)
	return storageMock
}
