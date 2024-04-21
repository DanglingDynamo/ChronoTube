package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

func Decrypt(text string) ([]byte, error) {
	bytes := []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
	block, err := aes.NewCipher([]byte(os.Getenv("ENCRYPTION_SECRET")))
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBDecrypter(block, bytes)

	decoded := make([]byte, len(cipherText))
	cfb.XORKeyStream(decoded, cipherText)
	return decoded, nil
}
