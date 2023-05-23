package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// GenerateRandom generates random string size N
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Encode encodes bytes array
func Encode(src []byte, key string) (string, error) {

	if len(key) == 0 {
		return "", errors.New("key is empty")
	}

	key32 := sha256.Sum256([]byte(key))

	aesblock, err := aes.NewCipher(key32[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce := key32[len(key32)-aesgcm.NonceSize():]

	data := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем

	based64 := base64.StdEncoding.EncodeToString(data)

	return based64, nil
}

// Decode decodes bytes array
func Decode(src string, key string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	key32 := sha256.Sum256([]byte(key))

	aesblock, err := aes.NewCipher(key32[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	// создаём вектор инициализации
	nonce := key32[len(key32)-aesgcm.NonceSize():]

	// расшифровываем
	decrypted, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	if len(decrypted) == 0 {
		decrypted = make([]byte, 0)

	}
	return decrypted, nil
}
