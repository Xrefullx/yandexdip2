package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCard_Encrypt_Decrypt(t *testing.T) {
	key := "secret_key"

	tests := []struct {
		name       string
		key        string
		data       []byte
		requireErr bool
	}{
		{
			name:       "empty key",
			data:       []byte("some_text_to_encode"),
			requireErr: true,
		},
		{
			name: "empty data",
			key:  key,
			data: make([]byte, 0),
		},
		{
			name: "short data",
			key:  key,
			data: []byte("123"),
		},
		{
			name: "large data 10mb",
			key:  key,
			data: func() []byte {
				d, err := GenerateRandom(1024 * 1024 * 10)
				require.NoError(t, err)
				return d
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			encoded, err := Encode(tt.data, tt.key)

			if tt.requireErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			decoded, err := Decode(encoded, tt.key)
			require.NoError(t, err)

			require.Equal(t, tt.data, decoded)
		})
	}
}
