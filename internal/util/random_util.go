package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(len int) (string, error) {
	b, err := GenerateRandomBytes(len)
	return base64.URLEncoding.EncodeToString(b), err
}
