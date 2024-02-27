package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString() (string, error) {
	b := make([]byte, 24) // 24 bytes will result in 32 characters after base64 encoding
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
